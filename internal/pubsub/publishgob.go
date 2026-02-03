package pubsub

import (
	"bytes"
	"encoding/gob"
	amqp "github.com/rabbitmq/amqp091-go"
	"context"
)

func PublishGob[T any](ch *amqp.Channel, exchange, key string, val T) error {
	// Encode the value using gob
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(val)
	if err != nil {
		return err
	}
	body := buf.Bytes()

	// Publish the message
	if err != nil {
		return err
	}
	return ch.PublishWithContext(
		context.Background(),
		exchange, // exchange
		key,      // routing key
		false,    // mandatory
		false,    // immediate
		amqp.Publishing{
			ContentType: "application/gob",
			Body:        body,
		},
	)
}