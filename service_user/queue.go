package main

import (
	"context"
	"errors"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type QueuePublisher struct {
	Channel *amqp.Channel
	Queues  map[string]amqp.Queue
}

func NewQueuePublisher(dialUrl string, queueName ...string) QueuePublisher {
	conn, err := amqp.Dial(dialUrl)
	if err != nil {
		log.Panic("failed to create connection to amqp:", err.Error())
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Panic("failed to open channel amqp:", err.Error())
	}

	q := map[string]amqp.Queue{}
	for _, name := range queueName {
		newQ := CreateQueue(ch, name)
		q[name] = newQ
	}

	return QueuePublisher{Channel: ch, Queues: q}
}

func CreateQueue(ch *amqp.Channel, name string) amqp.Queue {
	q, err := ch.QueueDeclare(
		name,  // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)

	if err != nil {
		log.Panic("failed to open channel amqp:", err.Error())
	}

	return q
}

func (qb QueuePublisher) SendMessage(ctx context.Context, qname string, message []byte) error {
	q := qb.Queues[qname]
	if q.Name == "" {
		return errors.New("queue was not initialized")
	}

	return qb.Channel.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        message,
		})
}
