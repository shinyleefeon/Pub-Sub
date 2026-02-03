package pubsub

import (
    "bytes"
    "encoding/gob"
    "log"
    "fmt"

    amqp "github.com/rabbitmq/amqp091-go"
)

func SubscribeGob[T any](
    conn *amqp.Connection,
    exchange,
    queueName,
    key string,
    queueType SimpleQueueType, // an enum to represent "durable" or "transient"
    handler func(T) AckType,
) error {

    newchan, q, err :=DeclareAndBind(conn, exchange, queueName, key, queueType)
    if err != nil {
        return err
    }
    newchan.Qos(10,0,false) // prefetch count of 10
    delivChan, err := newchan.Consume(q.Name, "", false, false, false, false, nil)
    if err != nil {
        return err
    }

    go func() {
        for d := range delivChan {
            var msg T
            buf := bytes.NewBuffer(d.Body)
            decoder := gob.NewDecoder(buf)
            if err := decoder.Decode(&msg); err != nil {
                log.Printf("failed to decode message: %v", err)
                continue
            }
            
            switch handler(msg) {
                case Ack:
                    d.Ack(false)
                    fmt.Printf("acknowledged message: %v\n", msg)
                case NackRequeue:
                    d.Nack(false, true)
                    fmt.Printf("requeued message: %v\n", msg)
                case NackDiscard:
                    d.Nack(false, false)
                    fmt.Printf("discarded message: %v\n", msg)
                }
        }
    }()

    return nil
}
