package main

import (
	"fmt"
	"os"
	signal "os/signal"
	"strconv"
	"time"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	username, err := gamelogic.ClientWelcome()
	gameState := gamelogic.NewGameState(username)
	connectionString := "amqp://guest:guest@localhost:5672/"
	conn, err := amqp.Dial(connectionString)
	if err != nil {
		fmt.Printf("Error during connection: %v\n", err)
		os.Exit(1)
	}
	chann, err := conn.Channel()
	if err != nil {
		fmt.Printf("Error during channel creation: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()
	err = pubsub.SubscribeJSON(conn, routing.ExchangePerilDirect, fmt.Sprintf("%s.%s", routing.PauseKey, username), routing.PauseKey, pubsub.SimpleQueueTypeTransient, handlerPause(gameState))
	if err != nil {
		fmt.Printf("Error during publishing: %v\n", err)
		os.Exit(1)
	}
	err = pubsub.SubscribeJSON(conn, routing.ExchangePerilTopic, fmt.Sprintf("%s.%s", routing.ArmyMovesPrefix, username), routing.ArmyMovesPrefix+".*", pubsub.SimpleQueueTypeTransient, handlerMove(chann, gameState))
	if err != nil {
		fmt.Printf("Error during publishing: %v\n", err)
		os.Exit(1)
	}
	err = pubsub.SubscribeJSON(conn, routing.ExchangePerilTopic, "war", routing.WarRecognitionsPrefix+".*", pubsub.SimpleQueueTypeDurable, handlerWar(chann, gameState))

	if err != nil {
		fmt.Printf("Error during publishing: %v\n", err)
		os.Exit(1)
	}
mainLoop:
	for true {
		words := gamelogic.GetInput()
		if len(words) == 0 {
			continue
		}
		switch words[0] {
		case "spawn":
			gameState.CommandSpawn(words)
		case "move":
			move, err := gameState.CommandMove(words)
			if err == nil {
				pubsub.PublishJSON(chann, routing.ExchangePerilTopic, fmt.Sprintf("%s.%s", routing.ArmyMovesPrefix, username), move)
				fmt.Println("Move published")
			}
		case "status":
			gameState.CommandStatus()
		case "help":
			gamelogic.PrintClientHelp()
		case "spam":
			if len(words) < 2 {
				continue
			}
			cuantity, err := strconv.ParseInt(words[1], 10, 32)
			if err != nil {
				continue
			}
			for i := 0; i < int(cuantity); i++ {
				log := routing.GameLog{
					CurrentTime: time.Now(),
					Message:     gamelogic.GetMaliciousLog(),
					Username:    gameState.GetUsername(),
				}
				publishGameLog(chann, log)
			}
			fmt.Printf("Sending %d spam messages\n", cuantity)
		case "quit":
			gamelogic.PrintQuit()
			break mainLoop
		default:
			fmt.Println("Command not found")
		}
	}

}
func handlerPause(gs *gamelogic.GameState) func(routing.PlayingState) pubsub.AckType {
	return func(ps routing.PlayingState) pubsub.AckType {
		gs.HandlePause(ps)
		defer fmt.Print("> ")

		return pubsub.Ack
	}
}
func handlerMove(channel *amqp.Channel, gs *gamelogic.GameState) func(gamelogic.ArmyMove) pubsub.AckType {
	return func(move gamelogic.ArmyMove) pubsub.AckType {
		outcome := gs.HandleMove(move)
		recOfWar := gamelogic.RecognitionOfWar{
			Attacker: move.Player,
			Defender: gs.GetPlayerSnap(),
		}
		if outcome == gamelogic.MoveOutComeSafe {
			return pubsub.Ack
		}
		if outcome == gamelogic.MoveOutcomeMakeWar {
			pubsub.PublishJSON(channel, routing.ExchangePerilTopic, fmt.Sprintf("%s.%s", routing.WarRecognitionsPrefix, move.Player.Username), recOfWar)
			// if err != nil {
			// 	return pubsub.NackRequeue
			// }
			return pubsub.Ack
		}
		defer fmt.Print("> ")
		return pubsub.NackDiscard
	}
}
func handlerWar(channel *amqp.Channel, gs *gamelogic.GameState) func(gamelogic.RecognitionOfWar) pubsub.AckType {
	defer fmt.Print("> ")
	return func(row gamelogic.RecognitionOfWar) pubsub.AckType {
		outcome, winner, loser := gs.HandleWar(row)
		gameLog := routing.GameLog{CurrentTime: time.Now(), Username: gs.GetUsername()}
		var ackType pubsub.AckType
		switch outcome {
		case gamelogic.WarOutcomeNotInvolved:
			ackType = pubsub.NackRequeue
		case gamelogic.WarOutcomeOpponentWon:
			gameLog.Message = fmt.Sprintf("%s won a war against %s", winner, loser)
			ackType = pubsub.Ack
		case gamelogic.WarOutcomeYouWon:
			gameLog.Message = fmt.Sprintf("%s won a war against %s", winner, loser)
			ackType = pubsub.Ack
		case gamelogic.WarOutcomeDraw:
			gameLog.Message = fmt.Sprintf("A war between %s and %s resulted in a draw", winner, loser)
			ackType = pubsub.Ack
		default:
			ackType = pubsub.NackDiscard
		}

		err := publishGameLog(channel, gameLog)
		if err != nil {
			return pubsub.NackRequeue
		}

		return ackType
	}
}

func publishGameLog(chann *amqp.Channel, log routing.GameLog) error {
	return pubsub.PublishGob(chann, routing.ExchangePerilTopic, fmt.Sprintf("%s.%s", routing.GameLogSlug, log.Username), log)
}
func waitForSignal() {
	c := make(chan os.Signal, 1)
	signal.Notify(c)
	<-c
	fmt.Println("Shutting down")
	os.Exit(0)
}
