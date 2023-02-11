package service

import (
	"context"
	"fmt"

	"github.com/AssylzhanZharzhanov/axxonsoft-test-service/internal/domain"
)

type service struct {
	eventService domain.EventService
	repository   domain.TaskRepository
}

// NewService - creates a new service
func NewService(
	eventService domain.EventService,
	repository domain.TaskRepository,
) domain.TaskService {
	service := newBasicService(eventService, repository)
	return service
}

func newBasicService(
	eventService domain.EventService,
	repository domain.TaskRepository,
) domain.TaskService {
	return &service{
		eventService: eventService,
		repository:   repository,
	}
}

func (s *service) CreateTask(ctx context.Context, task *domain.Task) (domain.TaskID, error) {
	if err := task.Validate(); err != nil {
		return 0, err
	}

	return s.repository.Create(ctx, task)
}

func (s *service) GetTask(ctx context.Context, taskID domain.TaskID) (*domain.Task, error) {
	if taskID <= 0 {
		return nil, fmt.Errorf("task id required")
	}

	return s.repository.Get(ctx, taskID)
}

func (s *service) ListTasks(ctx context.Context, criteria domain.TaskSearchCriteria) ([]*domain.Task, error) {
	return s.repository.List(ctx, criteria)
}
