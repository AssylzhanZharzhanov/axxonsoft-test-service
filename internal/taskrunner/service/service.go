package service

import (
	"context"
	"net/http"

	"github.com/AssylzhanZharzhanov/axxonsoft-test-service/internal/domain"

	"github.com/go-kit/log"
)

type service struct {
	taskRepository domain.TaskRepository
}

// NewTaskRunnerService - creates a new service
func NewTaskRunnerService(taskRepository domain.TaskRepository, logger log.Logger) domain.TaskRunnerService {
	var service domain.TaskRunnerService
	{
		service = newBasicService(taskRepository)
	}
	return service
}

func newBasicService(taskRepository domain.TaskRepository) domain.TaskRunnerService {
	return &service{
		taskRepository: taskRepository,
	}
}

func (s *service) RunTask(ctx context.Context, taskID domain.TaskID) error {

	task, err := s.taskRepository.Get(ctx, taskID)
	if err != nil {
		task.StatusID = domain.StatusError
	}

	err = s.DoRequest(ctx, task)
	if err != nil {
		task.StatusID = domain.StatusError
	}

	err = s.taskRepository.Update(ctx, task)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) DoRequest(ctx context.Context, task *domain.Task) error {
	httpRequest, err := http.NewRequestWithContext(ctx, task.Method, task.URL, nil)
	if err != nil {
		task.StatusID = domain.StatusError
	}
	// @TODO add headers
	//httpRequest.Header.Set("Content-Type", "application/json")
	//httpRequest.Header.Set("X-Api-Key", s.secretAPIKey)

	httpResponse, err := http.DefaultClient.Do(httpRequest)
	if err != nil {
		return err
	}
	defer httpResponse.Body.Close()

	task.StatusID = domain.StatusDone
	task.HTTPStatusCode = httpResponse.StatusCode
	task.ContentLength = httpResponse.ContentLength

	return nil
}
