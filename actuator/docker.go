package actuator

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/google/uuid"
)

var dockerClient *client.Client

func initDockerClient() *client.Client {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	dockerClient = cli
	return cli
}

func ImagePull(_image string) {
	reader, err := dockerClient.ImagePull(context.Background(), _image, image.PullOptions{})

	if err != nil {
		panic(err)
	}
	defer reader.Close()

	bodyBytes, err := io.ReadAll(reader)
	if err != nil {
		panic(err)
	}
	bodyString := string(bodyBytes)
	fmt.Println("image pull result ", bodyString)
}

func RunInDocker(image string, cmd []string, mounts []mount.Mount, timeLimit int) RunResult {
	start := time.Now()
	ret := RunResult{
		ExitCode: 0,
		Stdout:   "",
		Stderr:   "",
	}

	cgroup, err := NewCgroupWrap(uuid.New().String())
	if err != nil {
		panic(err)
	}
	defer cgroup.Delete()

	stopTimeout := 5
	// 这里假设所有操作都能在 (timeLimit+5)s 内完成
	// @TODO 每个语言允许配置编译耗时
	containerRunTimeout := time.Duration((timeLimit + stopTimeout) * int(time.Second))
	ctx, cancel := context.WithTimeout(context.Background(), containerRunTimeout)
	defer cancel()

	// 默认使用 1G 的内存限制
	size1G := int64(1 * 1024 * 1024 * 1024)

	createStart := time.Since(start)
	resp, err := dockerClient.ContainerCreate(ctx, &container.Config{
		NetworkDisabled: true,
		StopTimeout:     &stopTimeout,
		Image:           image,
		Cmd:             cmd,
	}, &container.HostConfig{
		Mounts: mounts,
		Resources: container.Resources{
			CgroupParent: cgroup.Path,
			Memory:       size1G,
			MemorySwap:   size1G,
			Ulimits: []*container.Ulimit{
				{
					Name: "cpu",
					// cpu load 较高时，如果进程因为这个 limit 被干掉，cgroup usage_usec 数据会稍微小 0.05s
					// 这里为了方便判断是否为 TLE，设置为 timeLimit+1
					Hard: int64(timeLimit + 1),
					Soft: int64(timeLimit),
				},
				{
					Name: "fsize",
					Hard: size1G,
					Soft: size1G,
				},
				{
					Name: "nproc",
					Hard: 1024,
					Soft: 1024,
				},
				{
					Name: "nice",
					Hard: -1,
					Soft: -1,
				},
			},
		},
	}, nil, nil, "")
	createEnd := time.Since(start)
	if err != nil {
		panic(err)
	}

	if err := dockerClient.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		panic(err)
	}
	startEnd := time.Since(start)

	// 超时时停止容器
	time.AfterFunc(containerRunTimeout, func() {
		dockerClient.ContainerStop(context.Background(), resp.ID, container.StopOptions{})
	})

	// 因为上面超时会暂停，所以这里不需要设置超时
	statusCh, errCh := dockerClient.ContainerWait(context.Background(), resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			panic(err)
		}
	case result := <-statusCh:
		{
			ret.ExitCode = int(result.StatusCode)
		}
	}
	waitEnd := time.Since(start)

	getLogCtx, cancelLogCtx := context.WithTimeout(context.Background(), containerRunTimeout)
	defer cancelLogCtx()
	out, err := dockerClient.ContainerLogs(getLogCtx, resp.ID, container.LogsOptions{ShowStdout: true, ShowStderr: true})
	if err != nil {
		panic(err)
	}
	defer out.Close()

	var stdout bytes.Buffer
	stdoutWriter := bufio.NewWriter(&stdout)

	var stderr bytes.Buffer
	stderrWriter := bufio.NewWriter(&stderr)

	stdcopy.StdCopy(stdoutWriter, stderrWriter, out)
	stdoutWriter.Flush()
	stderrWriter.Flush()

	ret.Stdout = stdout.String()
	ret.Stderr = stderr.String()

	logEnd := time.Since(start)

	stat, err := cgroup.Cgroup.Stat()
	if err != nil {
		panic(err)
	}

	ret.MemoryUsage = int(stat.Memory.GetMaxUsage())
	ret.TimeCost = int(stat.CPU.GetUsageUsec() / 1000)

	removeCtx, cancelRemoveCtx := context.WithTimeout(context.Background(), containerRunTimeout)
	defer cancelRemoveCtx()
	dockerClient.ContainerRemove(removeCtx, resp.ID, container.RemoveOptions{})

	removeEnd := time.Since(start)

	if ret.ExitCode != 0 {
		statJSON, err := json.Marshal(stat)
		if err != nil {
			panic(err)
		}

		if timeLimit*1000 > ret.TimeCost {
			// 其他未知情况
			fmt.Printf("cgroup.stat: %s\n", statJSON)
		} else {
			// 超时停止
			// fmt.Printf("TLE [%d] time limit: %d, time cost: %d\n", ret.ExitCode, timeLimit*1000, ret.TimeCost)
		}
	}

	fmt.Printf("createStart [%s]\n", createStart)
	fmt.Printf("createEnd [%s]\n", createEnd)
	fmt.Printf("startEnd [%s]\n", startEnd)
	fmt.Printf("waitEnd [%s]\n", waitEnd)
	fmt.Printf("logEnd [%s]\n", logEnd)
	fmt.Printf("removeEnd [%s]\n", removeEnd)
	fmt.Printf("UsageUsec [%d]\n", stat.CPU.GetUsageUsec())
	return ret
}
