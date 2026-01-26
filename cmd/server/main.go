package main

import (
    "fmt"
    amqp "github.com/rabbitmq/amqp091-go"
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
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, os.Kill)
	<-signalChan
	fmt.Println("Shutting down Peril server...")
}


