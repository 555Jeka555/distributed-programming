package ampq

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"time"
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
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	go func() {
		for m := range msgs {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			body, err := r.handler.Handle(ctx, m.Body)
			cancel()

			err = channel.PublishWithContext(ctx,
				"",        // exchange
				m.ReplyTo, // routing key
				false,     // mandatory
				false,     // immediate
				amqp.Publishing{
					ContentType:   "application/json",
					CorrelationId: m.CorrelationId,
					Body:          body,
				})
			failOnError(err, "Failed to publish a message")
			log.Println("Reader body", string(body))

			err = m.Ack(false)
		}
	}()

	return err
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
