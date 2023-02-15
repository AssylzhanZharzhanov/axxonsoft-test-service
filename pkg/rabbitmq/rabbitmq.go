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

func DeclareBinding(amqpChan *amqp.Channel, exchangeName string, exchangeKind string) (amqp.Queue, error) {
	var (
		queue amqp.Queue
	)

	if err := DeclareExchange(amqpChan, exchangeName, exchangeKind); err != nil {
		return queue, err
	}

	queue, err := DeclareQueue(amqpChan, exchangeName)
	if err != nil {
		return queue, err
	}

	if err := BindQueue(amqpChan, queue.Name, "", exchangeName); err != nil {
		return queue, err
	}

	return queue, nil
}

func DeclareExchange(amqpChan *amqp.Channel, exchangeName string, exchangeKind string) error {
	return amqpChan.ExchangeDeclare(
		exchangeName,
		exchangeKind,
		true,
		false,
		false,
		false,
		nil,
	)
}

func DeclareQueue(amqpChan *amqp.Channel, queueName string) (amqp.Queue, error) {
	return amqpChan.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
}

func BindQueue(amqpChan *amqp.Channel, queueName, key, exchangeName string) error {
	return amqpChan.QueueBind(queueName, key, exchangeName, false, nil)
}
