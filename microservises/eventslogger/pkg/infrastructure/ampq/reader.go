package ampq

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

func NewReader(queueName string, handler IntegrationEventHandler) Reader {
	return Reader{
		queueName: queueName,
		handler:   handler,
	}
}

type Reader struct {
	queueName string
	handler   IntegrationEventHandler
}

func (r *Reader) Connect(channel *amqp.Channel, routingKeys []string) error {
	err := channel.ExchangeDeclare(
		"logs",  // name
		"topic", // type
		true,    // durable
		false,   // auto-deleted
		false,   // internal
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		return err
	}

	q, err := channel.QueueDeclare(
		r.queueName, // name (теперь используем имя из параметра)
		false,       // durable
		false,       // delete when unused
		false,       // exclusive (false для возможности нескольких потребителей)
		false,       // no-wait
		nil,         // arguments
	)
	if err != nil {
		return err
	}

	for _, routingKey := range routingKeys {
		err = channel.QueueBind(
			q.Name,     // queue name
			routingKey, // routing key
			"logs",     // exchange
			false,
			nil,
		)
		if err != nil {
			return err
		}
	}

	msgs, err := channel.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return err
	}

	go func() {
		for d := range msgs {
			if err := r.handler.Handle(d.Body); err != nil {
				log.Printf("Error handling message: %v", err)
			}
		}
	}()

	return nil
}
