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
func NewConsumerService(taskRunnerService domain.TaskRunnerService, amqpChan *amqp.Channel, queueName string, consumer string, logger log.Logger) domain.Consumer {
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

	go func() {
		for {
			select {
			case <-ctx.Done():
				_ = s.logger.Log()
				return
			case msg, ok := <-messages:
				if !ok {
					_ = s.logger.Log()
					return
				}

				err = s.executeTask(ctx, msg)
				_ = s.logger.Log(err)
			}
		}
	}()

	return nil
}

func (s *service) executeTask(ctx context.Context, msg amqp.Delivery) error {

	var event domain.Event
	if err := json.Unmarshal(msg.Body, &event); err != nil {
		return msg.Reject(true)
	}

	err := s.taskRunnerService.RunTask(ctx, event.TaskID)
	if err != nil {
		return err
	}

	return nil
}
