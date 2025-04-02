package main

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"os"
	"server/pkg/app/handler"
	"server/pkg/infrastructure/ampq"
)

func main() {
	conn, err := amqp.Dial(os.Getenv("RABBITMQ_URL"))
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	handler := handler.NewHandler()
	integrationEventHandler := ampq.NewIntegrationEventHandler(handler)
	reader := ampq.NewReader("logs", integrationEventHandler)

	var forever chan struct{}

	err = reader.Connect(ch, []string{"rankcalculator.*", "valuator.*"})
	failOnError(err, "Failed to connect to queue")

	<-forever
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
