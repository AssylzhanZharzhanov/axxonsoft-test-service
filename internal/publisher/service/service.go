package service

import (
	"context"

	"github.com/AssylzhanZharzhanov/axxonsoft-test-service/internal/domain"

	"github.com/go-kit/log"
	amqp "github.com/rabbitmq/amqp091-go"
)

type service struct {
	amqpConn *amqp.Connection
	amqpChan *amqp.Channel
}

func NewService(amqpConn *amqp.Connection, amqpChan *amqp.Channel, logger log.Logger) domain.Publisher {
	var service domain.Publisher
	{
		service = newBasicService(amqpConn, amqpChan)
		service = loggingServiceMiddleware(logger)(service)
	}
	return service
}

func newBasicService(amqpConn *amqp.Connection, amqpChan *amqp.Channel) domain.Publisher {
	return &service{
		amqpConn: amqpConn,
		amqpChan: amqpChan,
	}
}

func (s service) Publish(ctx context.Context, event *domain.Event) error {
	if err := s.amqpChan.PublishWithContext(
		ctx,
		"",    // exchange
		"key", // key
		false,
		false,
		amqp.Publishing{},
	); err != nil {
		return err
	}

	return nil
}
