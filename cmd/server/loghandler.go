package main

import (
	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	"fmt"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
)



func LogHandler() func(routing.GameLog) pubsub.AckType{
	return func(gamelog routing.GameLog) pubsub.AckType{
		defer fmt.Printf(">")
		err := gamelogic.WriteLog(gamelog)
		if err != nil {
			fmt.Println("failed to write GameLogEntry:", err)
			return pubsub.NackRequeue
		}
		return pubsub.Ack
	}
}
