package ampq

import (
	"context"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func NewReader(
	queueName string,
	handler IntegrationEventHandler,
) Reader {
	return Reader{
		queueName: queueName,
		handler:   handler,
	}
}

type Reader struct {
	queueName string
	handler   IntegrationEventHandler
}

func (r *Reader) ConnectReadChannel(channel *amqp.Channel) error {
	q, err := channel.QueueDeclare(
		r.queueName, // name
		false,       // durable
		false,       // delete
		// when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := channel.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	go func() {
		for d := range msgs {
			err = r.handler.Handle(context.Background(), d.Body)
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	return err
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
