package main

import (
	"fmt"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
	"time"
)


func PublishGameLog(ch *amqp.Channel, gs *gamelogic.GameState, LogMsg string) error{
	err := pubsub.PublishGob(
			ch,
			routing.ExchangePerilTopic,
			routing.GameLogSlug + "." + gs.GetUsername(),
			routing.GameLog{
				Message: LogMsg,
				CurrentTime: time.Now(),
				Username: gs.GetUsername(),
			},
		)
		if err != nil {
			fmt.Println("failed to publish GameLogEntry:", err)
			return err
		}
	return nil
	}