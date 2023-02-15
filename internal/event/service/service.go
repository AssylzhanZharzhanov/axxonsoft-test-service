package service

import (
	"context"
	"fmt"

	"github.com/AssylzhanZharzhanov/axxonsoft-test-service/internal/domain"

	"github.com/go-kit/log"
)

type service struct {
	publisher domain.Publisher
}

func NewService(publisher domain.Publisher, logger log.Logger) domain.EventService {
	var service domain.EventService
	{
		service = newBasicService(publisher)
		service = loggingServiceMiddleware(logger)(service)
	}
	return service
}

func newBasicService(publisher domain.Publisher) domain.EventService {
	return &service{
		publisher: publisher,
	}
}

func (s *service) RegisterEvent(ctx context.Context, event *domain.Event) error {
	if event == nil {
		return fmt.Errorf("event required")
	}

	if err := event.Validate(); err != nil {
		return err
	}

	if err := s.publisher.Publish(ctx, event); err != nil {
		return err
	}

	return nil
}
