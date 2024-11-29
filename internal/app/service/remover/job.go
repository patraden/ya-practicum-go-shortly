package remover

import "github.com/patraden/ya-practicum-go-shortly/internal/app/dto"

type Job struct {
	ID    int
	Tasks []dto.UserSlug
}

type JobResult struct {
	ID  int
	Err error
}

type WorkerFunc func(j Job) JobResult
