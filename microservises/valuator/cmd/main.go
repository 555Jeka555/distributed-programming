package main

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"server/pkg/app"
	"server/pkg/infrastructure/storage"
	"server/pkg/infrastructure/transport"
)

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_URL"),
	})
	keyValueStorage := storage.NewKeyValue(rdb)
	valuatorService := app.NewValuatorService(keyValueStorage)
	ctx := context.Background()

	handler := transport.NewHandler(ctx, valuatorService)

	r := mux.NewRouter()

	r.HandleFunc("/", handler.Index).Methods("GET")
	r.HandleFunc("/summary", handler.Summary).Methods("POST")
	r.HandleFunc("/about", handler.About).Methods("GET")

	log.Println("Starting server on :8082")
	http.ListenAndServe(":8082", r)
}
