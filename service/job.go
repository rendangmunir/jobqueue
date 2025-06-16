package service

import (
	"context"
	"jobqueue/entity"
	"jobqueue/queue"
	"github.com/google/uuid"
	_interface "jobqueue/interface"
)

type jobService struct {
	jobRepo _interface.JobRepository
	jobQueue *queue.JobQueue
}

// Initiator ...
type Initiator func(s *jobService) *jobService

func (q jobService) GetAllJobs(ctx context.Context) ([]*entity.Job, error) {
	return q.jobRepo.FindAll(ctx)
}

func (q jobService) GetJobStatusSummary(ctx context.Context) (map[string]int32, error) {
	return q.jobRepo.CountByStatus(ctx)
}

func (q jobService) GetJobById(ctx context.Context, id string) (*entity.Job, error) {
	return q.jobRepo.FindByID(ctx, id)
}

func (q jobService) Enqueue(ctx context.Context, taskName string) (string, error) {
	id := uuid.NewString() 
	job := &entity.Job{
		ID:       id,
		Task:     taskName,
		Status:   "pending", 
		Attempts: 0,
	}

	if err := q.jobRepo.Save(ctx, job); err != nil {
		return "", err
	}

	// Add to JobQueue for processing
	q.jobQueue.Enqueue(job)
	
	return id, nil
}

func (q jobService) Update(ctx context.Context, job *entity.Job) error {
	return q.jobRepo.Save(ctx, job)
}

// NewJobService ...
func NewJobService() Initiator {
	return func(s *jobService) *jobService {
		return s
	}
}

// SetJobRepository ...
func (i Initiator) SetJobRepository(jobRepository _interface.JobRepository) Initiator {
	return func(s *jobService) *jobService {
		i(s).jobRepo = jobRepository
		return s
	}
}

func (i Initiator) SetJobQueue(jobQueue *queue.JobQueue) Initiator {
	return func(s *jobService) *jobService {
		i(s).jobQueue = jobQueue
		return s
	}
}

// Build ...
func (i Initiator) Build() _interface.JobService {
	return i(&jobService{})
}
