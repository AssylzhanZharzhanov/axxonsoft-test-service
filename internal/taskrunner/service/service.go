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

func NewTaskRunnerService(taskRepository domain.TaskRepository, logger log.Logger) domain.TaskRunnerService {
	var service domain.TaskRunnerService
	{
		service = newBasicService(taskRepository)
	}
	return service
}

func newBasicService(taskService domain.TaskRepository) domain.TaskRunnerService {
	return &service{
		taskRepository: taskService,
	}
}

func (s *service) RunTask(ctx context.Context, taskID domain.TaskID) error {

	task, err := s.taskRepository.Get(ctx, taskID)
	if err != nil {
		task.StatusID = domain.StatusError
	}

	url := task.URL
	method := task.Method

	httpRequest, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		task.StatusID = domain.StatusError
	}
	// @TODO add headers
	//httpRequest.Header.Set("Content-Type", "application/json")
	//httpRequest.Header.Set("X-Api-Key", s.secretAPIKey)

	httpResponse, err := http.DefaultClient.Do(httpRequest)
	if err != nil {
		task.StatusID = domain.StatusError
	}
	defer httpResponse.Body.Close()

	task.HTTPStatusCode = httpResponse.StatusCode
	task.ContentLength = httpResponse.ContentLength

	err = s.taskRepository.Update(ctx, task)
	if err != nil {
		return err
	}

	return nil
}
