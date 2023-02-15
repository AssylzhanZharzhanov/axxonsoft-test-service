package service

import (
	"context"
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"time"

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

	dataBytes, err := json.Marshal(&event)
	if err != nil {
		return err
	}

	msg := &amqp.Publishing{
		Timestamp: time.Now().UTC(),
		Body:      dataBytes,
	}

	if err := s.publisher.Publish(ctx, msg); err != nil {
		return err
	}

	return nil
}
