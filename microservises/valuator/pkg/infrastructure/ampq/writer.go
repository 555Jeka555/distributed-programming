package ampq

import (
	"context"
	"encoding/json"
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
		true,  // durable
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

func (w *writer) WriteExchange(evt event.Event) error {
	err := w.channel.ExchangeDeclare(
		"events", // name
		"topic",  // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = w.channel.PublishWithContext(ctx,
		"events",   // exchange
		evt.Type(), // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        serialize(evt),
		})
	failOnError(err, "Failed to publish a message")

	return nil
}

func serialize(evt event.Event) []byte {
	switch e := evt.(type) {
	case event.SimilarityCalculated:
		msg := SimilarityCalculatedMessage{
			Type:       e.Type(),
			TextID:     e.TextID,
			Similarity: e.Similarity,
		}

		msgSerialized, err := json.Marshal(msg)
		failOnError(err, "Failed to declare an exchange")

		return msgSerialized
	default:
		log.Printf("Unknown event type: %s", e.Type())
		return nil
	}
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

type SimilarityCalculatedMessage struct {
	Type       string `json:"type"`
	TextID     string `json:"text_id"`
	Similarity int    `json:"similarity"`
}
