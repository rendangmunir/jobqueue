package queue

import (
	"context"
	"fmt"
	"jobqueue/entity"
	_interface "jobqueue/interface"
	// "math/rand"
	"time"
)

// CreateJobHandler returns a handler func to process jobs
func CreateJobHandler(jobService _interface.JobService) func(ctx context.Context, job *entity.Job) {
    return func(ctx context.Context, job *entity.Job) {
        maxAttempts := 5
        for job.Attempts < int32(maxAttempts) {
            job.Attempts++
            job.Status = "running"
            _ = jobService.Update(ctx, job)
            fmt.Printf("Processing job %s attempt %d\n", job.ID, job.Attempts)

            success := process(job)

            if success {
                job.Status = "completed"
                _ = jobService.Update(ctx, job)
                fmt.Printf("Job %s completed\n", job.ID)
                return
            }

            job.Status = "failed"
            _ = jobService.Update(ctx, job)
            fmt.Printf("Job %s failed attempt %d\n", job.ID, job.Attempts)

            if job.Task == "unstable-job" && job.Attempts < int32(maxAttempts) {
                time.Sleep(500 * time.Millisecond)
                continue
            }

            return
        }
    }
}

func process(job *entity.Job) bool {
	time.Sleep(1500 * time.Millisecond) 

	if job.Task == "unstable-job" {
		if job.Attempts <= 2 {
			return false
		}
		return true
	}
	// uncomment for failure simulation
	// simulate 50% chance failure
	// if it fails, return false
	// return rand.Intn(2) == 0
	
	return true
}
