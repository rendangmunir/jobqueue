package _interface

import (
	"context"
	"jobqueue/entity"
)

type JobService interface {
	Enqueue(ctx context.Context, taskName string) (string, error)
	GetAllJobs(ctx context.Context) ([]*entity.Job, error)
	GetJobStatusSummary(ctx context.Context) (map[string]int32, error)
	GetJobById(ctx context.Context, id string) (*entity.Job, error)
	Update(ctx context.Context, job *entity.Job) error
}

type JobRepository interface {
	Save(ctx context.Context, job *entity.Job) error
	FindByID(ctx context.Context, id string) (*entity.Job, error)
	CountByStatus(ctx context.Context) (map[string]int32, error)
	FindAll(ctx context.Context) ([]*entity.Job, error)
}
