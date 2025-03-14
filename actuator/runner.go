package actuator

import (
	"fmt"
	"log"
	"math"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/Pigeon-Developer/pigeon-oj-judge/solution"
	"github.com/docker/docker/api/types/mount"
)

type Runner struct {
	job      *solution.JudgeJob
	image    string
	buildCmd string
	runCmd   string
	Result   int    // 对应 db 中的 result
	Info     string // CE/RE 时对应的错误信息

	basePath     string // 用户代码与编译产物都会放在这个目录下，是此程序可以访问到的目录结构
	hostBasePath string // 用户代码与编译产物都会放在这个目录下，是 docker host 可以访问到的目录结构

	// 存放每个测试点的详细信息
	checkPointDetail map[string]UserCodeRunResult
}

// 使用内置的语言镜像
func NewRunnerBuiltin(job *solution.JudgeJob, basePath, hostBasePath string) *Runner {
	runtimeConfig := RuntimeRegistry[job.Data.Language]
	return &Runner{
		job:      job,
		image:    runtimeConfig.Image,
		buildCmd: runtimeConfig.BuildCmd,
		runCmd:   runtimeConfig.RunCmd,
		Result:   0,
		Info:     "",

		basePath:     basePath,
		hostBasePath: hostBasePath,

		checkPointDetail: make(map[string]UserCodeRunResult, 0),
	}
}

func (r *Runner) Judge() {
	ret := r.Build()
	if ret.ExitCode != 0 {
		r.Result = solution.Result_CE
		r.Info = ret.Stderr
		return
	}

	r.Run()
}

func (r *Runner) Build() RunResult {
	solutionPath := path.Join(r.basePath, "solution", strconv.Itoa(r.job.Data.SolutionId))
	artifactPath := path.Join(solutionPath, "artifacts")

	hostSolutionPath := path.Join(r.hostBasePath, "solution", strconv.Itoa(r.job.Data.SolutionId))
	hostArtifactPath := path.Join(hostSolutionPath, "artifacts")

	os.MkdirAll(solutionPath, os.ModePerm)
	os.MkdirAll(artifactPath, os.ModePerm)
	writeFile(path.Join(solutionPath, "user_code"), r.job.Data.Code)

	buildTimeLimit := 5

	ret := RunInDocker(r.image, []string{"bash", "-c", r.buildCmd}, []mount.Mount{
		{
			Type:   mount.TypeBind,
			Source: hostSolutionPath,
			Target: "/mount/source-code",
		},
		{
			Type:   mount.TypeBind,
			Source: hostArtifactPath,
			Target: "/mount/artifacts",
		},
	}, buildTimeLimit)

	return ret
}

// 返回值为是否需要停止后续判题
func (r *Runner) RunAndCompare(checkPointName, hostArtifactPath, hostInDataPath, hostOutDataPath, outDataPath, systemOutDataPath string) bool {
	// artifactPath/inDataPath/outDataPath 需要使用 docker host 可以访问的路径
	// systemOutDataPath 是正确答案的路径，使用此程序可以访问的路径
	runTimeLimit := int(math.Ceil(r.job.Data.TimeLimit))

	runResult := RunInDocker(r.image, []string{"bash", "-c", r.runCmd}, []mount.Mount{
		{
			ReadOnly: false,
			Type:     mount.TypeBind,
			Source:   hostArtifactPath,
			Target:   "/mount/artifacts",
		},
		{
			ReadOnly: true,
			Type:     mount.TypeBind,
			Source:   hostInDataPath,
			Target:   "/app/data.in",
		},
		{
			ReadOnly: false,
			Type:     mount.TypeBind,
			Source:   hostOutDataPath,
			Target:   "/app/data.out",
		},
	}, runTimeLimit)

	r.checkPointDetail[checkPointName] = UserCodeRunResult{
		ExitCode:    runResult.ExitCode,
		MemoryUsage: runResult.MemoryUsage,
		TimeCost:    runResult.TimeCost,
		Stdout:      runResult.Stdout,
		Stderr:      runResult.Stderr,
		Match:       -1,
	}

	match := -1

	// 是否超时
	isTimeout := false
	if runResult.TimeCost >= (runTimeLimit*1000) || SIGXCPU == runResult.ExitCode {
		isTimeout = true
	}

	if isTimeout {
		r.Result = solution.Result_TLE
		return false
	}

	if runResult.ExitCode != 0 {
		fmt.Printf("运行失败 %d - %s - [%d] \n[%s]\n[%s]\n\n", r.job.Data.SolutionId, checkPointName, runResult.ExitCode, runResult.Stdout, runResult.Stderr)
		// 有一个异常就直接跳过后续的测试数据
		r.Result = solution.Result_RE
		return false
	}

	isMatch := CompareLineByLine(systemOutDataPath, outDataPath)

	if isMatch {
		match = 0
	} else {
		match = 1
	}

	r.checkPointDetail[checkPointName] = UserCodeRunResult{
		ExitCode:    runResult.ExitCode,
		MemoryUsage: runResult.MemoryUsage,
		TimeCost:    runResult.TimeCost,
		Stdout:      runResult.Stdout,
		Stderr:      runResult.Stderr,
		Match:       match,
	}

	if !isMatch {
		r.Result = solution.Result_WA

		return false
	}

	return true
}

func (r *Runner) Run() {
	hostSolutionPath := path.Join(r.hostBasePath, "solution", strconv.Itoa(r.job.Data.SolutionId))
	solutionPath := path.Join(r.basePath, "solution", strconv.Itoa(r.job.Data.SolutionId))

	dataPath := solution.GetSolutionDataPath(r.job.SourceID)
	problemDataPath := path.Join(dataPath, strconv.Itoa(r.job.Data.ProblemId))

	dataHostPath := solution.GetSolutionDataHostPath(r.job.SourceID)
	problemDataHostPath := path.Join(dataHostPath, strconv.Itoa(r.job.Data.ProblemId))

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

		outfileName := strings.Replace(e.Name(), ".in", ".out", 1)

		// 这里用于容器文件 bind，需要使用 host path
		hostInDataPath := path.Join(problemDataHostPath, e.Name())
		hostOutDataPath := path.Join(hostSolutionPath, "output", outfileName)
		hostArtifactPath := path.Join(hostSolutionPath, "artifacts")

		outDataPath := path.Join(solutionPath, "output", outfileName)
		systemOutput := path.Join(problemDataPath, e.Name())

		_, err := os.Create(hostOutDataPath)
		if err != nil {
			panic(err)
		}

		hasNext := r.RunAndCompare(e.Name(), hostArtifactPath, hostInDataPath, hostOutDataPath, outDataPath, systemOutput)
		if !hasNext {

			// @TODO 这里需要从测试点详情里面提取出 info
			return
		}
	}
}
