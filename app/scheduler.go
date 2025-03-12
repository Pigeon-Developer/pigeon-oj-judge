package app

import (
	"fmt"
	"sync"
	"time"

	"github.com/Pigeon-Developer/pigeon-oj-judge/actuator"
	"github.com/Pigeon-Developer/pigeon-oj-judge/solution"
)

// MaxConcurrentJudges defines the maximum number of concurrent judge operations
var MaxConcurrentJudges = 4

// WorkerPool manages a pool of workers for concurrent judgement
type WorkerPool struct {
	jobChan chan *solution.JudgeJob
	wg      sync.WaitGroup
}

// NewWorkerPool creates a new worker pool with n workers
func NewWorkerPool(n int) *WorkerPool {
	return &WorkerPool{
		jobChan: make(chan *solution.JudgeJob, n),
	}
}

// Start starts the worker pool
func (wp *WorkerPool) Start() {
	for i := range MaxConcurrentJudges {
		wp.wg.Add(1)
		go func(workerID int) {
			defer wp.wg.Done()
			for job := range wp.jobChan {
				result := actuator.JudgeUserSubmit(job)
				job.UpdateResult(result)
			}
		}(i)
	}
}

// Submit submits a job to the worker pool
func (wp *WorkerPool) Submit(job *solution.JudgeJob) {
	wp.jobChan <- job
}

// Close closes the worker pool
func (wp *WorkerPool) Close() {
	close(wp.jobChan)
	wp.wg.Wait()
}

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
	// Create a semaphore channel with capacity MaxConcurrentJudges
	semaphore := make(chan struct{}, MaxConcurrentJudges)

	for {
		// Try to acquire a semaphore slot
		select {
		case semaphore <- struct{}{}: // Successfully acquired a slot
			// Fetch and process a new job only when we have capacity
			go func() {
				defer func() {
					// Release the semaphore when the job is done
					<-semaphore
				}()

				job, err := fetchSolutionFromPool()
				if err == nil {
					result := actuator.JudgeUserSubmit(job)
					job.UpdateResult(result)
				} else {
					fmt.Println(err)
				}
			}()
		default:
			// Maximum concurrency reached, wait before checking again
			// No new tasks will be fetched until a worker becomes available
		}

		time.Sleep(500 * time.Millisecond)
	}
}
