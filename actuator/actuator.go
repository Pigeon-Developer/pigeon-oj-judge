package actuator

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/google/uuid"

	"github.com/Pigeon-Developer/pigeon-oj-judge/solution"
)

type RunResult struct {
	ExitCode    int
	Stdout      string
	Stderr      string
	MemoryUsage int // 单位 bytes
	TimeCost    int // 单位 ms
}

const BasePath = "/etc/pigeon-oj-judge"

func writeFile(filePath string, content string) {
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	file.WriteString(content)
}

func RunInDocker(image string, cmd []string, mounts []mount.Mount, timeLimit int) RunResult {
	ret := RunResult{
		ExitCode: 0,
		Stdout:   "",
		Stderr:   "",
	}

	cgroup, err := NewCgroupWrap(uuid.New().String())
	if err != nil {
		panic(err)
	}
	defer cgroup.Cgroup.Delete()

	stopTimeout := 5
	// 这里假设所有操作都能在 (timeLimit+5)s 内完成
	// @TODO 每个语言允许配置编译耗时
	buildTimeout := time.Duration((timeLimit + stopTimeout) * int(time.Second))
	ctx, cancel := context.WithTimeout(context.Background(), buildTimeout)
	defer cancel()

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	// 默认使用 1G 的内存限制
	size1G := int64(1 * 1024 * 1024 * 1024)

	resp, err := cli.ContainerCreate(ctx, &container.Config{
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
			// @TODO 这里应该均匀的分配核心
			CpusetCpus: "0",
			Ulimits: []*container.Ulimit{
				{
					Name: "rss",
					Hard: size1G,
					Soft: size1G,
				},
				{
					Name: "stack",
					Hard: size1G,
					Soft: size1G,
				},
				{
					Name: "data",
					Hard: size1G,
					Soft: size1G,
				},
				{
					Name: "fsize",
					Hard: size1G,
					Soft: size1G,
				},
				{
					Name: "nice",
					Hard: -1,
					Soft: -1,
				},
			},
		},
	}, nil, nil, "")
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		panic(err)
	}

	// 超时时停止容器
	time.AfterFunc(buildTimeout, func() {
		cli.ContainerStop(context.Background(), resp.ID, container.StopOptions{})
	})

	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
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

	out, err := cli.ContainerLogs(ctx, resp.ID, container.LogsOptions{ShowStdout: true, ShowStderr: true})
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

	stat, err := cgroup.Cgroup.Stat()
	if err != nil {
		panic(err)
	}

	statJSON, err := json.Marshal(stat)
	if err != nil {
		panic(err)
	} else {
		fmt.Printf("cgroup.stat: %s\n", statJSON)
	}

	ret.MemoryUsage = int(stat.Memory.GetMaxUsage())
	ret.TimeCost = int(stat.CPU.GetUsageUsec() / 1000)

	cli.ContainerRemove(ctx, resp.ID, container.RemoveOptions{})

	return ret
}

// 编译代码
func buildUserSubmitCode(job *solution.JudgeJob) RunResult {
	solutionPath := path.Join(BasePath, "solution", strconv.Itoa(job.Data.SolutionId))
	artifactPath := path.Join(solutionPath, "artifacts")

	os.MkdirAll(solutionPath, os.ModePerm)
	os.MkdirAll(artifactPath, os.ModePerm)
	writeFile(path.Join(solutionPath, "source_code"), job.Data.Code)

	buildTimeLimit := 5

	runtimeConfig := RuntimeRegistry[job.Data.Language]
	ret := RunInDocker(runtimeConfig.Image, []string{"bash", "-l", "-c", runtimeConfig.BuildCmd}, []mount.Mount{
		{
			Type:   mount.TypeBind,
			Source: solutionPath,
			Target: "/mount/source-code",
		},
		{
			Type:   mount.TypeBind,
			Source: artifactPath,
			Target: "/mount/artifacts",
		},
	}, buildTimeLimit)

	return ret
}

// 运行用户提交，使用测试数据中 in 得到用户的 out
func runUserSubmitCode(job *solution.JudgeJob) map[string]UserCodeRunResult {
	ret := make(map[string]UserCodeRunResult, 0)

	solutionPath := path.Join(BasePath, "solution", strconv.Itoa(job.Data.SolutionId))

	dataPath := solution.GetSolutionDataPath(job.SourceID)
	problemDataPath := path.Join(dataPath, strconv.Itoa(job.Data.ProblemId))

	os.MkdirAll(path.Join(solutionPath, "output"), os.ModePerm)

	entries, err := os.ReadDir(problemDataPath)
	if err != nil {
		log.Fatal(err)
	}

	for _, e := range entries {
		if e.IsDir() {
			continue
		}

		if !strings.HasSuffix(e.Name(), ".in") {
			continue
		}

		inDataPath := path.Join(problemDataPath, e.Name())
		outDataPath := path.Join(solutionPath, "output", strings.Replace(e.Name(), ".in", ".out", 1))
		artifactPath := path.Join(solutionPath, "artifacts")

		_, err := os.Create(outDataPath)
		if err != nil {
			panic(err)
		}

		runTimeLimit := int(math.Ceil(job.Data.TimeLimit))

		runtimeConfig := RuntimeRegistry[job.Data.Language]
		runResult := RunInDocker(runtimeConfig.Image, []string{"bash", "-l", "-c", runtimeConfig.RunCmd}, []mount.Mount{
			{
				ReadOnly: false,
				Type:     mount.TypeBind,
				Source:   artifactPath,
				Target:   "/mount/artifacts",
			},
			{
				ReadOnly: true,
				Type:     mount.TypeBind,
				Source:   inDataPath,
				Target:   "/app/data.in",
			},
			{
				ReadOnly: false,
				Type:     mount.TypeBind,
				Source:   outDataPath,
				Target:   "/app/data.out",
			},
		}, runTimeLimit)

		ret[e.Name()] = UserCodeRunResult{
			ExitCode:    runResult.ExitCode,
			MemoryUsage: runResult.MemoryUsage,
			TimeCost:    runResult.TimeCost,
			Stdout:      runResult.Stdout,
			Stderr:      runResult.Stderr,
			Match:       -1,
		}
	}

	return ret
}

// 判断用户输出是否与数据一致
func judgeUserSubmitCode(job *solution.JudgeJob, runResult map[string]UserCodeRunResult) int {
	for _, v := range runResult {
		if v.ExitCode != 0 {
			return solution.Result_RE
		}
	}

	solutionPath := path.Join(BasePath, "solution", strconv.Itoa(job.Data.SolutionId))
	solutionOutputPath := path.Join(solutionPath, "output")

	dataPath := solution.GetSolutionDataPath(job.SourceID)
	problemDataPath := path.Join(dataPath, strconv.Itoa(job.Data.ProblemId))

	entries, err := os.ReadDir(problemDataPath)
	if err != nil {
		log.Fatal(err)
	}

	// @TODO 结果改为 int，记录每个用例下具体的 result 状态
	result := make(map[string]bool, 32)
	isAllMatch := true

	for _, e := range entries {
		// 筛选出需要对比的文件
		if e.IsDir() {
			continue
		}

		if !strings.HasSuffix(e.Name(), ".out") {
			continue
		}

		//
		systemOutput := path.Join(problemDataPath, e.Name())
		userOutput := path.Join(solutionOutputPath, e.Name())

		systemFile, err := os.Open(systemOutput)
		if err != nil {
			log.Fatal(err)
		}
		defer systemFile.Close()

		userFile, err := os.Open(userOutput)
		if err != nil {
			log.Fatal(err)
		}
		defer userFile.Close()

		systemScanner := bufio.NewScanner(systemFile)
		userScanner := bufio.NewScanner(userFile)

		isMatch := true
		// optionally, resize scanner's capacity for lines over 64K, see next example
		for systemScanner.Scan() {
			userScanner.Scan()

			if systemScanner.Text() != userScanner.Text() {
				isMatch = false
				break
			}
		}

		result[e.Name()] = isMatch
		isAllMatch = isAllMatch && isMatch

		if err := systemScanner.Err(); err != nil {
			log.Fatal(err)
		}
	}

	fmt.Printf("用户的答案是否正确 %v \n", result)

	if isAllMatch {
		return solution.Result_AC
	} else {
		return solution.Result_WA
	}
}

func JudgeUserSubmit(job *solution.JudgeJob) int {
	job.UpdateResult(solution.Result_CI)
	compileResult := buildUserSubmitCode(job)
	if compileResult.ExitCode != 0 {
		return solution.Result_CE
	}
	job.UpdateResult(solution.Result_RJ)
	runResult := runUserSubmitCode(job)
	return judgeUserSubmitCode(job, runResult)
}
