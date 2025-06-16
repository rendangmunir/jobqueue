package queue

import (
	"context"
	"jobqueue/entity"
	"sync"
)

type JobQueue struct {
	queue chan *entity.Job
	wg    sync.WaitGroup
}


func NewJobQueue(bufferSize int) *JobQueue {
	return &JobQueue{
		queue: make(chan *entity.Job, bufferSize),
	}
}

func (q *JobQueue) Enqueue(job *entity.Job) {
	q.queue <- job
}

func (q *JobQueue) StartWorkers(ctx context.Context, n int, handler func(ctx context.Context, job *entity.Job)) {
	for i := 0; i < n; i++ {
		q.wg.Add(1)
		go func(workerID int) {
			defer q.wg.Done()
			for {
				select {
				case job := <-q.queue:
					handler(ctx, job)
				case <-ctx.Done():
					return
				}
			}
		}(i)
	}
}

// Wait blocks until all workers exit
func (q *JobQueue) Wait() {
	q.wg.Wait()
}
