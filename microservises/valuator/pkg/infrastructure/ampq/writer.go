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

func (w *writer) Write(body []byte) ([]byte, error) {
	q, err := w.channel.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := w.channel.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack факт получения = факт доставки
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

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

	var res []byte
	for d := range msgs {
		if corrId == d.CorrelationId {
			res = d.Body
			failOnError(err, "Failed to convert body to integer")
			break
		}
	}

	return res, err
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
