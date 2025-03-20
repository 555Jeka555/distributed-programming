package ampq

import (
	"context"
	"github.com/gofrs/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"server/pkg/app/event"
	"time"
)

func NewWriter(
	queueName string,
	channel *amqp.Channel,
) event.Writer {
	return &writer{
		queueName: queueName,
		channel:   channel,
	}
}

type writer struct {
	queueName string
	channel   *amqp.Channel
}

func (w *writer) Write(body []byte) error {
	q, err := w.channel.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	corrId := uuid.Must(uuid.NewV4()).String()

	err = w.channel.PublishWithContext(ctx,
		"",          // exchange
		w.queueName, // routing key
		false,       // mandatory
		false,       // immediate
		amqp.Publishing{
			CorrelationId: corrId,
			ContentType:   "application/json",
			ReplyTo:       q.Name,
			Body:          body,
		})
	failOnError(err, "Failed to publish a message")
	log.Println("Writer body", string(body))

	return err
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
