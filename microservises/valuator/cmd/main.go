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
	"server/pkg/infrastructure/redis/repo"
	"server/pkg/infrastructure/transport"
	"time"
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

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body := "Hello World!"
	err = ch.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	failOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s\n", body)

	textRepo := repo.NewTextRepository(rdb)
	valuatorService := service.NewValuatorService(textRepo)
	textQueryService := query.NewTextQueryService(textRepo)

	handler := transport.NewHandler(ctx, valuatorService, textQueryService)

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
