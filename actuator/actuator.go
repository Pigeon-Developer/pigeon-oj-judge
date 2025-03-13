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
	defer cgroup.Delete()

	stopTimeout := 5
	// 这里假设所有操作都能在 (timeLimit+5)s 内完成
	// @TODO 每个语言允许配置编译耗时
	containerRunTimeout := time.Duration((timeLimit + stopTimeout) * int(time.Second))
	ctx, cancel := context.WithTimeout(context.Background(), containerRunTimeout)
	defer cancel()

	// 默认使用 1G 的内存限制
	size1G := int64(1 * 1024 * 1024 * 1024)

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
	if err != nil {
		panic(err)
	}

	if err := dockerClient.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		panic(err)
	}

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

	stat, err := cgroup.Cgroup.Stat()
	if err != nil {
		panic(err)
	}

	ret.MemoryUsage = int(stat.Memory.GetMaxUsage())
	ret.TimeCost = int(stat.CPU.GetUsageUsec() / 1000)

	removeCtx, cancelRemoveCtx := context.WithTimeout(context.Background(), containerRunTimeout)
	defer cancelRemoveCtx()
	dockerClient.ContainerRemove(removeCtx, resp.ID, container.RemoveOptions{})

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

	return ret
}

// 编译代码
func buildUserSubmitCode(job *solution.JudgeJob) RunResult {
	solutionPath := path.Join(BasePath, "solution", strconv.Itoa(job.Data.SolutionId))
	artifactPath := path.Join(solutionPath, "artifacts")

	os.MkdirAll(solutionPath, os.ModePerm)
	os.MkdirAll(artifactPath, os.ModePerm)
	writeFile(path.Join(solutionPath, "user_code"), job.Data.Code)

	buildTimeLimit := 5

	runtimeConfig := RuntimeRegistry[job.Data.Language]
	ret := RunInDocker(runtimeConfig.Image, []string{"bash", "-c", runtimeConfig.BuildCmd}, []mount.Mount{
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

	dataHostPath := solution.GetSolutionDataHostPath(job.SourceID)
	problemDataHostPath := path.Join(dataHostPath, strconv.Itoa(job.Data.ProblemId))

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

		// 这里用于容器文件 bind，需要使用 host path
		inDataPath := path.Join(problemDataHostPath, e.Name())
		outDataPath := path.Join(solutionPath, "output", strings.Replace(e.Name(), ".in", ".out", 1))
		artifactPath := path.Join(solutionPath, "artifacts")

		_, err := os.Create(outDataPath)
		if err != nil {
			panic(err)
		}

		runTimeLimit := int(math.Ceil(job.Data.TimeLimit))

		runtimeConfig := RuntimeRegistry[job.Data.Language]
		runResult := RunInDocker(runtimeConfig.Image, []string{"bash", "-c", runtimeConfig.RunCmd}, []mount.Mount{
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

		if runResult.ExitCode != 0 {
			if runResult.TimeCost >= (runTimeLimit * 1000) {
				// 这里是超时
			} else {
				fmt.Printf("运行失败 %d - %s - [%d] \n[%s]\n[%s]\n\n", job.Data.SolutionId, e.Name(), runResult.ExitCode, runResult.Stdout, runResult.Stderr)
			}

			// 有一个异常就直接跳过后续的测试数据
			// @TODO 应该允许每个 job 单独配置
			return ret
		}
	}

	return ret
}

func CompareLineByLine(file1, file2 string) bool {
	f1, err := os.Open(file1)
	if err != nil {
		log.Fatal(err)
	}
	defer f1.Close()

	f2, err := os.Open(file2)
	if err != nil {
		log.Fatal(err)
	}
	defer f2.Close()

	scanner1 := bufio.NewScanner(f1)
	scanner2 := bufio.NewScanner(f2)

	for {
		leftHasData := scanner1.Scan()
		rightHasData := scanner2.Scan()

		if leftHasData && rightHasData {
			l := scanner1.Text()
			r := scanner2.Text()
			if l != r {
				return false
			} else {
				continue
			}
		}

		if !leftHasData && !rightHasData {
			return true
		}

		// @TODO 这里需要看看剩余的字符是否全为 \n
		return false
	}
}

// 判断用户输出是否与数据一致
func judgeUserSubmitCode(job *solution.JudgeJob, runResult map[string]UserCodeRunResult) int {
	for _, v := range runResult {
		// 是否超时
		timeLimit := int(math.Ceil(job.Data.TimeLimit)) * 1000
		if SIGXCPU == v.ExitCode || v.TimeCost >= timeLimit {
			// @TODO 这里应该细分一下内存超出限制的情况
			return solution.Result_TLE
		}
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

		systemOutput := path.Join(problemDataPath, e.Name())
		userOutput := path.Join(solutionOutputPath, e.Name())

		isMatch := CompareLineByLine(systemOutput, userOutput)

		result[e.Name()] = isMatch
		isAllMatch = isAllMatch && isMatch
	}

	if !isAllMatch {
		fmt.Printf("用户答案 %d %v \n", job.Data.SolutionId, result)
	}

	if isAllMatch {
		return solution.Result_AC
	} else {
		return solution.Result_WA
	}
}

func JudgeUserSubmit(job *solution.JudgeJob) int {
	compileResult := buildUserSubmitCode(job)
	if compileResult.ExitCode != 0 {
		fmt.Printf("编译失败 %d \n[%s]\n[%s]\n\n", job.Data.SolutionId, compileResult.Stdout, compileResult.Stderr)
		return solution.Result_CE
	}
	job.UpdateResult(solution.Result_RJ)
	runResult := runUserSubmitCode(job)
	return judgeUserSubmitCode(job, runResult)
}
