package ampq

import (
	"context"
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"server/pkg/app/event"
	"time"
)

func NewWriter(
	channel *amqp.Channel,
) event.Writer {
	return &writer{
		channel: channel,
	}
}

type writer struct {
	channel *amqp.Channel
}

func (w *writer) WriteExchange(evt event.Event) error {
	err := w.channel.ExchangeDeclare(
		"logs",  // name
		"topic", // type
		true,    // durable
		false,   // auto-deleted
		false,   // internal
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = w.channel.PublishWithContext(ctx,
		"logs",     // exchange
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
	case event.RankCalculated:
		msg := RankCalculatedMessage{
			Type:   e.Type(),
			TextID: e.TextID,
			Rank:   e.Rank,
		}

		msgSerialized, err := json.Marshal(msg)
		failOnError(err, "Failed to declare an exchange")

		return msgSerialized
	default:
		log.Printf("Unknown event type: %s", e.Type())
		return nil
	}
}

type RankCalculatedMessage struct {
	Type   string  `json:"type"`
	TextID string  `json:"text_id"`
	Rank   float64 `json:"rank"`
}
