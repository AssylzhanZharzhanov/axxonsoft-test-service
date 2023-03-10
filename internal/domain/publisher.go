package domain

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Publisher - represents publisher behaviour
type Publisher interface {
	// Publish - publishes event in message queue
	//
	Publish(ctx context.Context, msg *amqp.Publishing) error
}
