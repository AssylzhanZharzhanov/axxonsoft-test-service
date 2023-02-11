package service

import (
	"context"
	"fmt"

	"github.com/AssylzhanZharzhanov/axxonsoft-test-service/internal/domain"
)

const (
	cacheDuration = 3600
)

type service struct {
	eventService    domain.EventService
	repository      domain.TaskRepository
	redisRepository domain.TaskRedisRepository
}

// NewService - creates a new service
func NewService(
	eventService domain.EventService,
	repository domain.TaskRepository,
	redisRepository domain.TaskRedisRepository,
) domain.TaskService {
	service := newBasicService(eventService, repository, redisRepository)
	return service
}

func newBasicService(
	eventService domain.EventService,
	repository domain.TaskRepository,
	redisRepository domain.TaskRedisRepository,
) domain.TaskService {
	return &service{
		eventService:    eventService,
		repository:      repository,
		redisRepository: redisRepository,
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

	// Get from cache
	key := fmt.Sprintf("?task/%d", taskID)
	cachedTask, err := s.redisRepository.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	if cachedTask != nil {
		return cachedTask, nil
	}

	// Get from storage
	result, err := s.repository.Get(ctx, taskID)
	if err != nil {
		return nil, err
	}

	// Set new value to cache
	err = s.redisRepository.Set(ctx, key, result, cacheDuration)
	if err != nil {
		return nil, err
	}

	return result, err
}

func (s *service) ListTasks(ctx context.Context, criteria domain.TaskSearchCriteria) ([]*domain.Task, error) {
	return s.repository.List(ctx, criteria)
}
