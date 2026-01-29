package pubsub

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

type SimpleQueueType int

const (
	Durable SimpleQueueType = iota
	Transient
)

func DeclareAndBind(
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType, // an enum to represent "durable" or "transient"
) (*amqp.Channel, amqp.Queue, error) {
	newChan, err := conn.Channel()
	if err != nil {
		return nil, amqp.Queue{}, err
	}

	
	durableparam := false
	autoDelete := false
	exclusive := false

	switch queueType {
	case Durable:
		durableparam = true
	case Transient:
		autoDelete = true
		exclusive = true
	}



	q, err := newChan.QueueDeclare(
		queueName,
		durableparam,
		autoDelete,
		exclusive,
		false,
		nil,
	)
	if err != nil {
		return nil, amqp.Queue{}, err
	}

	err = newChan.QueueBind(
		q.Name,
		key,
		exchange,
		false,
		nil,
	)
	if err != nil {
		return nil, amqp.Queue{}, err
	}

	return newChan, q, nil
}