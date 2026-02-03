package main

import (
	"fmt"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	amqp "github.com/rabbitmq/amqp091-go"

	
)

func handlerWar(gs *gamelogic.GameState, ch *amqp.Channel) func(msgBody gamelogic.RecognitionOfWar) pubsub.AckType {
	return func( msgBody gamelogic.RecognitionOfWar) pubsub.AckType {
		defer fmt.Print("> ")
		warOutcome, winner, loser := gs.HandleWar(msgBody)
		LogMsg := ""
		pubBool := false
	
	if warOutcome == gamelogic.WarOutcomeYouWon || warOutcome == gamelogic.WarOutcomeOpponentWon {
		LogMsg = fmt.Sprintf("%s won a war against %s", winner, loser)
		pubBool = true
	} else if warOutcome == gamelogic.WarOutcomeDraw {
		LogMsg = fmt.Sprintf("A war between %s and %s resulted in a draw", msgBody.Attacker.Username, msgBody.Defender.Username)
		pubBool = true
	}

	if pubBool {
		err := PublishGameLog(ch, gs, LogMsg)
		if err != nil {
			fmt.Println("failed to publish GameLogEntry:", err)
			return pubsub.NackRequeue
		}
	}


	switch warOutcome {
	case gamelogic.WarOutcomeNotInvolved:
		return pubsub.NackRequeue
	case gamelogic.WarOutcomeOpponentWon:
		return pubsub.Ack
	case gamelogic.WarOutcomeYouWon:
		return pubsub.Ack
	case gamelogic.WarOutcomeDraw:
		return pubsub.Ack
	case gamelogic.WarOutcomeNoUnits:
		return pubsub.NackDiscard
	default:
		fmt.Println("Unknown war outcome:", warOutcome)
		return pubsub.NackDiscard
		}
	}
}

	
