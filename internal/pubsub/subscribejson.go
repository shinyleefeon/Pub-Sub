package pubsub

import (
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)


func SubscribeJSON[T any](
    conn *amqp.Connection,
    exchange,
    queueName,
    key string,
    queueType SimpleQueueType, // an enum to represent "durable" or "transient"
    handler func(T),
) error {

	newchan, q, err :=DeclareAndBind(conn, exchange, queueName, key, queueType)
	if err != nil {
		return err
	}

	delivChan, err := newchan.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	go func() {
		for d := range delivChan {
			var msg T
			if err := json.Unmarshal(d.Body, &msg); err != nil {
				log.Printf("failed to unmarshal message: %v", err)
				continue
			}
			handler(msg)
			d.Ack(false)
		}
	}()

	return nil
}

