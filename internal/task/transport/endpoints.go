package transport

import (
	"context"

	"github.com/AssylzhanZharzhanov/axxonsoft-test-service/internal/domain"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/log"
)

type Endpoints struct {
	CreateTaskEndpoint endpoint.Endpoint
	GetTaskEndpoint    endpoint.Endpoint
	ListTasksEndpoint  endpoint.Endpoint
}

func NewEndpoints(service domain.TaskService, log log.Logger) Endpoints {
	return Endpoints{
		CreateTaskEndpoint: MakeCreateTaskEndpoint(service),
		GetTaskEndpoint:    MakeGetTaskEndpoint(service),
		ListTasksEndpoint:  MakeListTasksEndpoint(service),
	}
}

type createTaskRequest struct {
	Task *domain.Task
}

type createTaskResponse struct {
	TaskID domain.TaskID
	Err    error
}

func MakeCreateTaskEndpoint(service domain.TaskService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(createTaskRequest)

		resp, err := service.CreateTask(ctx, req.Task)
		return createTaskResponse{
			TaskID: resp,
			Err:    err,
		}, nil
	}
}

func MakeGetTaskEndpoint(service domain.TaskService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getTaskRequest)

		resp, err := service.GetTask(ctx, req.TaskID)
		return getTaskResponse{
			Task: resp,
			Err:  err,
		}, nil
	}
}

type getTaskRequest struct {
	TaskID domain.TaskID
}

type getTaskResponse struct {
	Task *domain.Task
	Err  error
}

func MakeListTasksEndpoint(service domain.TaskService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(listTaskRequest)

		resp, total, err := service.ListTasks(ctx, req.Criteria)
		return listTaskResponse{
			Tasks: resp,
			Total: int(total),
			Err:   err,
		}, nil
	}
}

type listTaskRequest struct {
	Criteria domain.TaskSearchCriteria
}

type listTaskResponse struct {
	Tasks []*domain.Task
	Total int
	Err   error
}
