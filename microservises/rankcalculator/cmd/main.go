package main

import (
	"github.com/go-redis/redis/v8"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"os"
	"server/pkg/app/handler"
	"server/pkg/app/service"
	"server/pkg/infrastructure/ampq"
	"server/pkg/infrastructure/centrifugo"
	"server/pkg/infrastructure/redis/provider"
	"server/pkg/infrastructure/redis/repo"
)

func main() {
	mainRdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_MAIN_URL"),
	})

	shards := map[string]*redis.Client{
		"RU":   redis.NewClient(&redis.Options{Addr: os.Getenv("REDIS_RU_URL")}),
		"EU":   redis.NewClient(&redis.Options{Addr: os.Getenv("REDIS_EU_URL")}),
		"ASIA": redis.NewClient(&redis.Options{Addr: os.Getenv("REDIS_ASIA_URL")}),
	}

	conn, err := amqp.Dial(os.Getenv("RABBITMQ_URL"))
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	writer := ampq.NewWriter(ch)

	shardManager := repo.NewShardManager(
		mainRdb,
		shards,
		map[string]string{
			"RU": "RU",
			"FR": "EU",
			"DE": "EU",
			"AE": "ASIA",
			"IN": "ASIA",
		},
	)
	rankCalculatorRepo := repo.NewShardTextRepository(shardManager)
	rankCalculatorService := service.NewRankCalculatorService(rankCalculatorRepo, writer)

	centrifugoClient := centrifugo.NewCentrifugoClient()
	textProvider := provider.NewTextProvider(rankCalculatorRepo)
	handler := handler.NewHandler(rankCalculatorService, textProvider, centrifugoClient)
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
