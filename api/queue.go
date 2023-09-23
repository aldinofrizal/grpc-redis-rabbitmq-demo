package main

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type QueueConsumer struct {
	Channel *amqp.Channel
}

func NewQueueConsumer(dialUrl string) QueueConsumer {
	conn, err := amqp.Dial(dialUrl)
	if err != nil {
		log.Panic("failed to create connection to amqp:", err.Error())
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Panic("failed to open channel amqp:", err.Error())
	}

	return QueueConsumer{Channel: ch}
}

func (qc QueueConsumer) AddConsumer(name string, handler func(amqp.Delivery) error) {
	q, err := qc.Channel.QueueDeclare(
		name,  // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)

	if err != nil {
		log.Panic("failed to declare queue:", err.Error())
	}

	msgs, err := qc.Channel.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)

	if err != nil {
		log.Panic("failed to register consumer:", err.Error())
	}

	go func() {
		for d := range msgs {
			log.Printf("received a %s message: %v", q.Name, string(d.Body))
			handler(d)
		}
	}()

}
