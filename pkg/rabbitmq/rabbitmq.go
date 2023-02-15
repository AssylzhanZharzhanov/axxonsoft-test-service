package rabbitmq

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

func NewRabbitMQConnection(uri string) (*amqp.Connection, error) {
	conn, err := amqp.Dial(uri)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
