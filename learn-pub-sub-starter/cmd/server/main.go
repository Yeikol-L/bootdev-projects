package main

import (
	"fmt"
	"os"
	signal "os/signal"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril server...")
	gamelogic.PrintServerHelp()
	connectionString := "amqp://guest:guest@localhost:5672/"
	conn, err := amqp.Dial(connectionString)
	if err != nil {
		fmt.Printf("Error during connection: %v", err)
		os.Exit(1)
	}
	defer conn.Close()
	chann, err := conn.Channel()
	if err != nil {
		fmt.Printf("Error during publishing: %v\n", err)
		os.Exit(1)
	}
	defer chann.Close()
	logChann, _, err := pubsub.DeclareAndBind(conn, routing.ExchangePerilTopic, routing.GameLogSlug, routing.GameLogSlug+".*", pubsub.SimpleQueueTypeDurable)
	if err != nil {
		fmt.Printf("Error during publishing: %v\n", err)
		os.Exit(1)
	}
	err = pubsub.SuscribeGob(conn, routing.ExchangePerilTopic, routing.GameLogSlug, routing.GameLogSlug+".*", pubsub.SimpleQueueTypeDurable, func(log routing.GameLog) pubsub.AckType {
		gamelogic.WriteLog(log)
		defer fmt.Println(log)
		return pubsub.Ack
	})
	if err != nil {
		fmt.Printf("Error during suscribing: %v\n", err)
		os.Exit(1)
	}
	defer logChann.Close()

mainLoop:
	for true {
		input := gamelogic.GetInput()
		if len(input) == 0 {
			continue
		}
		switch input[0] {
		case "pause":
			fmt.Println("Sending pause message")
			pubsub.PublishJSON(chann, routing.ExchangePerilDirect, routing.PauseKey,
				routing.PlayingState{
					IsPaused: true,
				},
			)
		case "resume":
			fmt.Println("Sending resume message")
			pubsub.PublishJSON(chann, routing.ExchangePerilDirect, routing.PauseKey,
				routing.PlayingState{
					IsPaused: false,
				},
			)
		case "quit":
			fmt.Println("Exiting app")
			break mainLoop
		default:
			fmt.Println("Command not found")
		}
	}
	os.Exit(0)
}

func waitForSignal() {
	c := make(chan os.Signal, 1)
	signal.Notify(c)
	<-c
	fmt.Println("Shutting down")
	os.Exit(0)
}
