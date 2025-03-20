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
	"server/pkg/infrastructure/ampq"
	"server/pkg/infrastructure/redis/repo"
	"server/pkg/infrastructure/transport"
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

	textRepo := repo.NewTextRepository(rdb)
	textQueryService := query.NewTextQueryService(textRepo)

	ctx := context.Background()
	writer := ampq.NewWriter("text", ch)
	handler := transport.NewHandler(ctx, writer, textQueryService)

	r := mux.NewRouter()

	r.HandleFunc("/", handler.Index).Methods("GET")
	r.HandleFunc("/summary", handler.Summary).Methods("POST")
	r.HandleFunc("/about", handler.About).Methods("GET")

	log.Println(fmt.Sprintf("Starting server on %s", os.Getenv("LISTENING_SERVER_PORT")))
	http.ListenAndServe(os.Getenv("LISTENING_SERVER_PORT"), r)
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
