package app

import (
	"fmt"
	"time"

	"github.com/Pigeon-Developer/pigeon-oj-judge/actuator"
	"github.com/Pigeon-Developer/pigeon-oj-judge/solution"
)

func fetchSolutionFromPool(jobChan chan<- *solution.JudgeJob) error {
	var soluton *solution.Solution
	var err error
	for _, instance := range solution.InstancePool {
		soluton, err = instance.Source.GetOne()
		if err == nil {
			jobChan <- &solution.JudgeJob{
				SourceID: instance.ID,
				Data:     soluton,
			}
			return nil
		}
	}

	return err
}

func RunLoop(maxConcurrent int, emptyWait int) {
	// Create job channel with buffer to avoid blocking
	jobChan := make(chan *solution.JudgeJob, maxConcurrent)

	// Create a semaphore to track active judge operations
	sem := make(chan struct{}, maxConcurrent)

	// Start job consumer goroutines
	for range maxConcurrent {
		go func() {
			for job := range jobChan {
				sem <- struct{}{} // Acquire semaphore
				result := actuator.JudgeUserSubmit(job)
				job.UpdateResult(result)
				<-sem // Release semaphore
			}
		}()
	}

	start := time.Now()

	for {
		if len(sem) < maxConcurrent {
			err := fetchSolutionFromPool(jobChan)
			if err == nil {
				continue
			}

			fmt.Println("fetchSolutionFromPool ", err)
			elapsed := time.Since(start)
			elapsed = elapsed.Round(time.Millisecond)
			fmt.Printf("程序已经运行了: %s\n", elapsed)
		}

		// 获取任务存在报错，或者任务队列已满
		time.Sleep(time.Duration(emptyWait) * time.Second)
	}
}
