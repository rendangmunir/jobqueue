package inmemrepo

import (
	"context"
	"errors"
	"jobqueue/entity"
	_interface "jobqueue/interface"
	"sync"
)

type jobRepository struct {
	mu      sync.RWMutex
	inMemDb map[string]*entity.Job
}

// Save Job
func (t *jobRepository) Save(ctx context.Context, job *entity.Job) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.inMemDb[job.ID] = job
	return nil
}

// Find Job By ID
func (t *jobRepository) FindByID(ctx context.Context, id string) (*entity.Job, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	job, exists := t.inMemDb[id]
	if !exists {
		return nil, errors.New("job not found")
	}
	return job, nil
}

// FindAll Job
func (t *jobRepository) FindAll(ctx context.Context) ([]*entity.Job, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	var jobs []*entity.Job
	for _, job := range t.inMemDb {
		jobs = append(jobs, job)
	}
	return jobs, nil
}

func (t *jobRepository) CountByStatus(ctx context.Context) (map[string]int32, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	counts := map[string]int32{
		"pending":   0,
		"running":   0,
		"failed":    0,
		"completed": 0,
	}

	for _, job := range t.inMemDb {
		if _, ok := counts[job.Status]; ok {
			counts[job.Status]++
		}
	}

	return counts, nil
}

// Initiator ...
type Initiator func(s *jobRepository) *jobRepository

// NewJobRepository ...
func NewJobRepository() Initiator {
	return func(q *jobRepository) *jobRepository {
		return q
	}
}

// SetInMemConnection set database client connection
func (i Initiator) SetInMemConnection(inMemDb map[string]*entity.Job) Initiator {
	return func(s *jobRepository) *jobRepository {
		i(s).inMemDb = inMemDb
		return s
	}
}

// Build ...
func (i Initiator) Build() _interface.JobRepository {
	return i(&jobRepository{})
}
