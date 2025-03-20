package main

import (
	"github.com/go-redis/redis/v8"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"os"
	"server/pkg/app/event"
	"server/pkg/app/service"
	"server/pkg/infrastructure/ampq"
	"server/pkg/infrastructure/redis/repo"
)

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_URL"),
	})

	conn, err := amqp.Dial(os.Getenv("RABBITMQ_URL"))
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	rankCalculatorRepo := repo.NewTextRepository(rdb)
	rankCalculatorService := service.NewRankCalculatorService(rankCalculatorRepo)

	handler := event.NewHandler(rankCalculatorService)
	integrationEventHandler := ampq.NewIntegrationEventHandler(handler)
	reader := ampq.NewReader("text", integrationEventHandler)

	var forever chan struct{}

	err = reader.ConnectReadChannel(ch)
	failOnError(err, "Failed to connect to ReadChannel")

	<-forever
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
