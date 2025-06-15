package query

import (
	"context"
	_dataloader "jobqueue/delivery/graphql/dataloader"
	"jobqueue/delivery/graphql/resolver"
	"jobqueue/entity"
	_interface "jobqueue/interface"
)

type JobQuery struct {
	jobService _interface.JobService
	dataloader *_dataloader.GeneralDataloader
}

func (q JobQuery) Jobs(ctx context.Context) ([]resolver.JobResolver, error) {
	jobs, err := q.jobService.GetAllJobs(ctx)
	if err != nil {
		return nil, err
	}

	resolvers := make([]resolver.JobResolver, len(jobs))
	for i, job := range jobs {
		resolvers[i] = resolver.JobResolver{
			Data:       *job,
			JobService: q.jobService,
			Dataloader: q.dataloader,
		}
	}
	return resolvers, nil
}

func (q JobQuery) Job(ctx context.Context, args struct {
	ID string
}) (*resolver.JobResolver, error) {
	// Fetch the job from service
	job, err := q.jobService.GetJobById(ctx, args.ID)
	if err != nil {
		return nil, err
	}

	return &resolver.JobResolver{
		Data:       *job,
		JobService: q.jobService,
		Dataloader: q.dataloader,
	}, nil
}

func (q JobQuery) JobStatus(ctx context.Context) (resolver.JobStatusResolver, error) {
	statusCounts, err := q.jobService.GetJobStatusSummary(ctx)
	if err != nil {
		return resolver.JobStatusResolver{}, err
	}

	data := entity.JobStatus{
		Pending:   statusCounts["pending"],
		Running:   statusCounts["running"],
		Failed:    statusCounts["failed"],
		Completed: statusCounts["completed"],
	}

	return resolver.JobStatusResolver{
		Data:       data,
		JobService: q.jobService,
		Dataloader: q.dataloader,
	}, nil
}

func NewJobQuery(jobService _interface.JobService,
	dataloader *_dataloader.GeneralDataloader) JobQuery {
	return JobQuery{
		jobService: jobService,
		dataloader: dataloader,
	}
}
