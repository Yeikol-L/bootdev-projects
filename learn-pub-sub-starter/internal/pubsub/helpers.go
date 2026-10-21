package pubsub

import (
	"bytes"
	"context"
	"encoding/gob"
	"encoding/json/v2"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type SimpleQueueType string

const (
	SimpleQueueTypeDurable   SimpleQueueType = "durable"
	SimpleQueueTypeTransient SimpleQueueType = "transient"
)

type AckType int

const (
	Ack AckType = iota
	NackRequeue
	NackDiscard
)

func PublishJSON[T any](ch *amqp.Channel, exchange, key string, val T) error {
	body, err := json.Marshal(val)
	if err != nil {
		return err
	}
	msg := amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	}
	return ch.PublishWithContext(context.Background(), exchange, key, false, false, msg)
}

func DeclareAndBind(
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType,
) (*amqp.Channel, amqp.Queue, error) {
	channel, err := conn.Channel()
	if err != nil {
		return nil, amqp.Queue{}, err
	}
	options := amqp.Table{
		"x-dead-letter-exchange": "peril_dlx",
	}
	queue, err := channel.QueueDeclare(queueName, queueType == SimpleQueueTypeDurable, queueType == SimpleQueueTypeTransient, queueType == SimpleQueueTypeTransient, false, options)
	if err != nil {
		return nil, amqp.Queue{}, err
	}
	err = channel.QueueBind(queue.Name, key, exchange, false, nil)
	return channel, queue, err
}

func SubscribeJSON[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType,
	handler func(T) AckType,
) error {
	chann, _, err := DeclareAndBind(conn, exchange, queueName, key, queueType)
	if err != nil {
		return err
	}
	delivery, err := chann.Consume(queueName, "", false, false, false, false, nil)
	if err != nil {
		return err
	}
	go func() {
		for v := range delivery {
			var body T
			err := json.Unmarshal(v.Body, &body)
			if err != nil {
				continue
			}
			result := handler(body)
			switch result {
			case Ack:
				v.Ack(false)
				fmt.Println("Ack")
			case NackRequeue:
				fmt.Println("NackRequeue")
				v.Nack(false, true)
			case NackDiscard:
				fmt.Println("NackDiscard")
				v.Nack(false, false)
			default:
				v.Ack(false)
			}
		}
	}()
	return nil
}

func PublishGob[T any](ch *amqp.Channel, exchange, key string, val T) error {
	var b bytes.Buffer
	err := gob.NewEncoder(&b).Encode(&val)
	if err != nil {
		return err
	}
	msg := amqp.Publishing{
		ContentType: "application/gob",
		Body:        b.Bytes(),
	}
	return ch.PublishWithContext(context.Background(), exchange, key, false, false, msg)
}

func SuscribeGob[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType,
	handler func(T) AckType,
) error {
	chann, _, err := DeclareAndBind(conn, exchange, queueName, key, queueType)
	if err != nil {
		return err
	}
	err = chann.Qos(10, 0, false)
	if err != nil {
		return err
	}
	delivery, err := chann.Consume(queueName, "", false, false, false, false, nil)
	if err != nil {
		return err
	}
	go func() {
		for v := range delivery {
			var body T
			err := gob.NewDecoder(bytes.NewBuffer(v.Body)).Decode(&body)
			if err != nil {
				continue
			}
			result := handler(body)
			switch result {
			case Ack:
				v.Ack(false)
				fmt.Println("Ack")
			case NackRequeue:
				fmt.Println("NackRequeue")
				v.Nack(false, true)
			case NackDiscard:
				fmt.Println("NackDiscard")
				v.Nack(false, false)
			default:
				v.Ack(false)
			}
		}
	}()
	return nil
}
