package app

import (
	"fmt"
	"time"

	"github.com/Pigeon-Developer/pigeon-oj-judge/actuator"
	"github.com/Pigeon-Developer/pigeon-oj-judge/solution"
)

func fetchSolutionFromPool() (*solution.JudgeJob, error) {
	var soluton *solution.Solution
	var err error
	for _, instance := range solution.InstancePool {
		soluton, err = instance.Source.GetOne()
		if err == nil {
			return &solution.JudgeJob{
				SourceID: instance.ID,
				Data:     soluton,
			}, nil
		}
	}

	return nil, err
}

func RunLoop() {
	for {
		job, err := fetchSolutionFromPool()
		if err == nil {
			// 需要判题
			result := actuator.JudgeUserSubmit(job)

			job.UpdateResult(result)
		} else {
			fmt.Println(err)
		}
		time.Sleep(4 * time.Second)
	}
}
