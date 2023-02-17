package service

import (
	"context"
	"encoding/json"
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
		service = loggingServiceMiddleware(logger)(service)
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

	err = s.MakeRequest(ctx, task)
	if err != nil {
		task.StatusID = domain.StatusError
	}

	err = s.taskRepository.Update(ctx, task)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) MakeRequest(ctx context.Context, task *domain.Task) error {
	httpRequest, err := http.NewRequestWithContext(ctx, task.Method, task.URL, nil)
	if err != nil {
		task.StatusID = domain.StatusError
	}

	headers, err := s.encodeHeaders(task.Headers)
	if err != nil {
		return err
	}

	for key, value := range headers {
		httpRequest.Header.Set(key, value)
	}

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

func (s *service) encodeHeaders(headers domain.JSONB) (map[string]string, error) {
	var result map[string]string
	err := json.Unmarshal(headers, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
