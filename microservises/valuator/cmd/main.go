package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"net/http"
	"os"
	"server/pkg/app/query"
	"server/pkg/app/service"
	"server/pkg/infrastructure/ampq"
	"server/pkg/infrastructure/redis/repo"
	"server/pkg/infrastructure/transport"
	"time"
)

func main() {
	mainRdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_MAIN_URL"),
		Password: os.Getenv("REDIS_MAIN_PASSWORD"),
		Username: "default",
	})

	pingCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := mainRdb.Ping(pingCtx).Result(); err != nil {
		log.Fatalf("Failed to connect to main Redis: %v", err)
	}
	log.Println("Successfully connected to main Redis")
	shards := map[string]*redis.Client{
		"RU": redis.NewClient(&redis.Options{
			Addr:     os.Getenv("REDIS_RU_URL"),
			Password: os.Getenv("REDIS_RU_PASSWORD"),
			Username: "default",
		}),
		"EU": redis.NewClient(&redis.Options{
			Addr:     os.Getenv("REDIS_EU_URL"),
			Password: os.Getenv("REDIS_EU_PASSWORD"),
			Username: "default",
		}),
		"ASIA": redis.NewClient(&redis.Options{
			Addr:     os.Getenv("REDIS_ASIA_URL"),
			Password: os.Getenv("REDIS_ASIA_PASSWORD"),
			Username: "default",
		}),
	}

	for region, client := range shards {
		if _, err := client.Ping(pingCtx).Result(); err != nil {
			log.Fatalf("Failed to connect to %s Redis shard: %v", region, err)
		}
		log.Printf("Successfully connected to %s Redis shard", region)
	}

	conn, err := amqp.Dial(os.Getenv("RABBITMQ_URL"))
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	regions := map[string]string{
		"RU": "RU",
		"FR": "EU",
		"DE": "EU",
		"AE": "ASIA",
		"IN": "ASIA",
	}

	shardManager := repo.NewShardManager(
		mainRdb,
		shards,
		regions,
	)
	textRepo := repo.NewShardTextRepository(shardManager)
	textQueryService := query.NewTextQueryService(textRepo)
	textService := *service.NewTextService(textRepo)

	ctx := context.Background()
	writer := ampq.NewWriter("text", ch)
	handler := transport.NewHandler(ctx, os.Getenv("JWT_KEY"), writer, textService, textQueryService, regions)

	r := mux.NewRouter()

	r.HandleFunc("/", handler.Index).Methods("GET")
	r.HandleFunc("/valuator/summary/create", handler.SummaryCreate).Methods("POST")
	r.HandleFunc("/valuator/summary", handler.Summary).Methods("GET")
	r.HandleFunc("/valuator/about", handler.About).Methods("GET")

	log.Println(fmt.Sprintf("Starting server on %s", os.Getenv("LISTENING_SERVER_PORT")))
	http.ListenAndServe(os.Getenv("LISTENING_SERVER_PORT"), r)
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
