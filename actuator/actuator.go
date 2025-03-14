package actuator

import (
	"bufio"
	"log"
	"os"

	"github.com/Pigeon-Developer/pigeon-oj-judge/solution"
)

type RunResult struct {
	ExitCode    int
	Stdout      string
	Stderr      string
	MemoryUsage int // 单位 bytes
	TimeCost    int // 单位 ms
}

func writeFile(filePath string, content string) {
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	file.WriteString(content)
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

func JudgeUserSubmit(job *solution.JudgeJob) int {
	runner := NewRunnerBuiltin(job, solution.GetBasePath(job.SourceID), solution.GetHostBasePath(job.SourceID))
	runner.Judge()
	return runner.Result
}
