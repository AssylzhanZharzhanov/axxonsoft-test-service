package service

import (
	"context"
	"encoding/json"
	"github.com/AssylzhanZharzhanov/axxonsoft-test-service/internal/domain"

	"github.com/go-kit/log"
	amqp "github.com/rabbitmq/amqp091-go"
)

type service struct {
	taskRunnerService domain.TaskRunnerService
	amqpChan          *amqp.Channel
	queueName         string
	consumer          string
	logger            log.Logger
}

// NewConsumerService creates a new consumer service with necessary dependencies.
func NewConsumerService(
	taskRunnerService domain.TaskRunnerService,
	amqpChan *amqp.Channel,
	queueName string,
	consumer string,
	logger log.Logger,
) domain.Consumer {
	return &service{
		taskRunnerService: taskRunnerService,
		amqpChan:          amqpChan,
		queueName:         queueName,
		consumer:          consumer,
		logger:            logger,
	}
}

func (s *service) Consume(ctx context.Context) error {
	messages, err := s.amqpChan.Consume(
		s.queueName,
		s.consumer,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	go s.worker(ctx, messages)

	chanErr := <-s.amqpChan.NotifyClose(make(chan *amqp.Error))
	return chanErr
}

func (s *service) worker(ctx context.Context, messages <-chan amqp.Delivery) {

	for msg := range messages {
		var event domain.Event
		if err := json.Unmarshal(msg.Body, &event); err != nil {
			if err = msg.Reject(true); err != nil {
				_ = s.logger.Log("error: %s", err)
			}
		}

		err := s.taskRunnerService.RunTask(ctx, event.TaskID)
		if err != nil {
			if err = msg.Reject(true); err != nil {
				_ = s.logger.Log("error: %s", err)
			}
		}

		err = msg.Ack(false)
		if err != nil {
			_ = s.logger.Log("error: %s", err)
		}
	}
}
