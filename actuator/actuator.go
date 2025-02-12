package actuator

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"math"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"

	"github.com/Pigeon-Developer/pigeon-oj-judge/solution"
)

const BasePath = "/etc/pigeon-oj-judge"

func writeFile(filePath string, content string) {
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	file.WriteString(content)
}

// 编译代码
func buildUserSubmitCode(job *solution.JudgeJob) {
	solutionPath := path.Join(BasePath, "solution", strconv.Itoa(job.Data.SolutionId))

	os.MkdirAll(solutionPath, os.ModePerm)
	writeFile(path.Join(solutionPath, "main.c"), job.Data.Code)

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	buildTimeLimit := 5

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		NetworkDisabled: true,
		StopTimeout:     &buildTimeLimit,
		Image:           "silkeh/clang:19-bookworm",
		Cmd:             []string{"clang", "/app/main.c", "-o", "/app/main.bin"},
	}, &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: solutionPath,
				Target: "/app",
			},
		},
	}, nil, nil, "")
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		panic(err)
	}

	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			panic(err)
		}
	case <-statusCh:
	}

	out, err := cli.ContainerLogs(ctx, resp.ID, container.LogsOptions{ShowStdout: true, ShowStderr: true})
	if err != nil {
		panic(err)
	}
	defer out.Close()

	stdcopy.StdCopy(os.Stdout, os.Stderr, out)

	cli.ContainerRemove(ctx, resp.ID, container.RemoveOptions{})
}

// 运行用户提交，使用测试数据中 in 得到用户的 out
func runUserSubmitCode(job *solution.JudgeJob) {
	solutionPath := path.Join(BasePath, "solution", strconv.Itoa(job.Data.SolutionId))
	problemDataPath := path.Join(BasePath, "data", strconv.Itoa(job.Data.ProblemId))

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
		exePath := path.Join(solutionPath, "main.bin")

		_, err := os.Create(outDataPath)
		if err != nil {
			panic(err)
		}

		ctx := context.Background()
		cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		if err != nil {
			panic(err)
		}

		runTimeLimit := int(math.Ceil(job.Data.TimeLimit))

		resp, err := cli.ContainerCreate(ctx, &container.Config{
			NetworkDisabled: true,
			StopTimeout:     &runTimeLimit,
			Image:           "silkeh/clang:19-bookworm",
			Cmd:             []string{"bash", "-l", "-c", " cat /app/data.in | /app/main.bin > /app/data.out"},
		}, &container.HostConfig{
			Mounts: []mount.Mount{
				{
					ReadOnly: false,
					Type:     mount.TypeBind,
					Source:   exePath,
					Target:   "/app/main.bin",
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
			},
		}, nil, nil, "")
		if err != nil {
			panic(err)
		}

		if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
			panic(err)
		}

		statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
		select {
		case err := <-errCh:
			if err != nil {
				panic(err)
			}
		case <-statusCh:
		}

		out, err := cli.ContainerLogs(ctx, resp.ID, container.LogsOptions{ShowStdout: true})
		if err != nil {
			panic(err)
		}
		defer out.Close()

		stdcopy.StdCopy(os.Stdout, os.Stderr, out)

		cli.ContainerRemove(ctx, resp.ID, container.RemoveOptions{})
	}
}

// 判断用户输出是否与数据一致
func judgeUserSubmitCode(job *solution.JudgeJob) int {
	solutionPath := path.Join(BasePath, "solution", strconv.Itoa(job.Data.SolutionId))
	solutionOutputPath := path.Join(solutionPath, "output")
	problemDataPath := path.Join(BasePath, "data", strconv.Itoa(job.Data.ProblemId))

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
	buildUserSubmitCode(job)
	runUserSubmitCode(job)
	return judgeUserSubmitCode(job)
}
