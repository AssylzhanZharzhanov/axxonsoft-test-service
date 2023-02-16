package service

import (
	"context"
	"github.com/AssylzhanZharzhanov/axxonsoft-test-service/internal/domain"

	"github.com/go-kit/log"
	amqp "github.com/rabbitmq/amqp091-go"
)

type service struct {
	amqpConn     *amqp.Connection
	amqpChan     *amqp.Channel
	exchangeName string
	exchangeKind string
}

func NewService(amqpConn *amqp.Connection, amqpChan *amqp.Channel, exchangeName string, exchangeKind string, logger log.Logger) domain.Publisher {
	var service domain.Publisher
	{
		service = newBasicService(amqpConn, amqpChan, exchangeName, exchangeKind)
		service = loggingServiceMiddleware(logger)(service)
	}
	return service
}

func newBasicService(amqpConn *amqp.Connection, amqpChan *amqp.Channel, exchangeName string, exchangeKind string) domain.Publisher {
	return &service{
		amqpConn:     amqpConn,
		amqpChan:     amqpChan,
		exchangeName: exchangeName,
		exchangeKind: exchangeKind,
	}
}

func (s *service) Publish(ctx context.Context, msg *amqp.Publishing) error {

	if err := s.amqpChan.PublishWithContext(
		ctx,
		s.exchangeName,
		"", // key
		false,
		false,
		*msg,
	); err != nil {
		return err
	}

	return nil
}
