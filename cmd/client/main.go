package main

import (
	"fmt"
	"os"
	"os/signal"
    amqp "github.com/rabbitmq/amqp091-go"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
)

const RabbitMQURL = "amqp://guest:guest@localhost:5672/"

func main() {
	fmt.Println("Starting Peril client...")
	a, err := amqp.Dial(RabbitMQURL)
	if err != nil {
		fmt.Println("Failed to connect to RabbitMQ:", err)
		return
	}
	defer a.Close()
	fmt.Println("Connected to RabbitMQ successfully")

	username, err := gamelogic.ClientWelcome()
	if err != nil {
		fmt.Println("Failed to get username:", err)
		return
	}
	fmt.Println("Welcome,", username)

	/*
	
	exchange := "peril_direct"
	queueName := "pause." + username
	routingKey := "pause"
	queueType := pubsub.Transient


	ch, q, err := pubsub.DeclareAndBind(a, exchange, queueName, routingKey, queueType)
	if err != nil {
		fmt.Println("Failed to declare and bind queue:", err)
		return
	}
	defer ch.Close()
	fmt.Println("Queue declared and bound successfully:", q.Name) */

	newState := gamelogic.NewGameState(username)
	pubsub.SubscribeJSON(a, routing.ExchangePerilDirect, "pause."+username, routing.PauseKey, pubsub.Transient, handlerPause(newState))


	for {
		command := gamelogic.GetInput()
		if len(command) == 0 {
			continue
		}
		switch command[0] {

		case "spawn":
			err := newState.CommandSpawn(command)
			if err != nil {
				fmt.Println("Error:", err)
			}
		
		case "move":
			_, err := newState.CommandMove(command)
			if err != nil {
				fmt.Println("Error:", err)
			}
		
		case "status":
			newState.CommandStatus()
		case "help":
			gamelogic.PrintClientHelp()
		case "spam":
			fmt.Println("Spamming not allowed yet!")
		case "quit":
			fmt.Println("Quitting the game...")
			return
		default:
			fmt.Println("Unknown command. Type 'help' for a list of commands.")
		}
	}

	

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, os.Kill)
	<-signalChan
	fmt.Println("Shutting down Peril client...")
}
