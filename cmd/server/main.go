package main

import (
	"os"
	"os/signal"
    "fmt"
    amqp "github.com/rabbitmq/amqp091-go"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
)

const RabbitMQURL = "amqp://guest:guest@localhost:5672/"

func main() {
	fmt.Println("Starting Peril server...")
	a, err := amqp.Dial(RabbitMQURL)
	if err != nil {
		fmt.Println("Failed to connect to RabbitMQ:", err)
		return
	}
	defer a.Close()
	fmt.Println("Connected to RabbitMQ successfully")
	publishCh, err := a.Channel()
	if err != nil {
		fmt.Println("Failed to open a channel:", err)
		return
	}
	defer publishCh.Close()
	fmt.Println("Channel opened successfully")

	err = pubsub.PublishJSON(publishCh, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{IsPaused: true})
	if err != nil {
		fmt.Println("Failed to publish message:", err)
		return
	}
	fmt.Println("Message published successfully")

	
	


	/*ch, q, err := pubsub.DeclareAndBind(a, routing.ExchangePerilTopic, routing.GameLogSlug, (routing.GameLogSlug + ".*"), pubsub.Durable)
	if err != nil {
		fmt.Println("Failed to declare and bind queue:", err)
		return
	}
	defer ch.Close()
	fmt.Println("Queue declared and bound successfully:", q.Name) */

	err = pubsub.SubscribeGob(a, routing.ExchangePerilTopic, routing.GameLogSlug, routing.GameLogSlug+".*", pubsub.Durable, LogHandler())
	if err != nil {
		fmt.Println("Failed to subscribe to GameLog messages:", err)
		return
	}
	fmt.Println("Subscribed to GameLog messages successfully")


	gamelogic.PrintServerHelp()

	for {
		command := gamelogic.GetInput()
		if len(command) == 0 {
			continue
		}
		switch command[0] {
			case "pause":
				err = pubsub.PublishJSON(publishCh, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{IsPaused: true})
				if err != nil {
					fmt.Println("Failed to publish pause message:", err)
				} else {
					fmt.Println("Pause message published successfully")
				}
			case "resume":
				err = pubsub.PublishJSON(publishCh, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{IsPaused: false})
				if err != nil {
					fmt.Println("Failed to publish resume message:", err)
				} else {
					fmt.Println("Resume message published successfully")
				}
			case "quit":
				fmt.Println("Quitting Peril server...")
				return
			default:
				fmt.Println("Unknown command:", command[0])
		}
	}
	

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, os.Kill)
	<-signalChan
	fmt.Println("Shutting down Peril server...")
}


