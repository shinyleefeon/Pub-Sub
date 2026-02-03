package main

import (
	"fmt"
	"os"
	"os/signal"
    amqp "github.com/rabbitmq/amqp091-go"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	"strconv"
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
	/*moveCh, q, err := pubsub.DeclareAndBind(a, routing.ExchangePerilTopic, "army_moves."+username, "army_moves.*", pubsub.Transient )
	if err != nil {
		fmt.Println("Failed to declare and bind queue for army moves:", err)
		return
	}
	defer moveCh.Close()*/
	
	moveCh, err := a.Channel()
	if err != nil {
		fmt.Println("Failed to create channel for publishing move:", err)
		return
	}
	defer moveCh.Close()
	warCh, err := a.Channel()
	if err != nil {
		fmt.Println("Failed to create channel for publishing war recognition:", err)
		return
	}
	defer warCh.Close()

	pubsub.SubscribeJSON(a, routing.ExchangePerilTopic, "army_moves."+username, "army_moves.*", pubsub.Transient, handlerMove(newState, warCh))
	fmt.Println("Queue declared and bound successfully for army moves")
	

	pubsub.SubscribeJSON(a, routing.ExchangePerilTopic, "war", routing.WarRecognitionsPrefix+".*", pubsub.Durable, handlerWar(newState, warCh))
	

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
			newMove, err := newState.CommandMove(command)
			if err != nil {
				fmt.Println("Error:", err)
			}
			
			pubsub.PublishJSON(moveCh, routing.ExchangePerilTopic, "army_moves."+username, newMove)
			fmt.Println("Published move")
		
		case "status":
			newState.CommandStatus()
		case "help":
			gamelogic.PrintClientHelp()
		case "spam":
			if len (command) < 2 {
				fmt.Println("Usage: spam <number_of_messages>")
				continue
			}
			numMessages, err := strconv.Atoi(command[1])
			if err != nil {
				fmt.Println("Invalid number of messages:", err)
				continue
			}
			for i := 0; i < numMessages; i++ {
				malLog := gamelogic.GetMaliciousLog()
				PublishGameLog(moveCh, newState, malLog)
			}
			fmt.Printf("Published %d malicious log messages\n", numMessages)
			continue
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
