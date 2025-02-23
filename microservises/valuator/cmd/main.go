package main

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"server/pkg/app/query"
	"server/pkg/app/service"
	"server/pkg/infrastructure/redis/repo"
	"server/pkg/infrastructure/transport"
)

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_URL"),
	})

	textRepo := repo.NewTextRepository(rdb)
	valuatorService := service.NewValuatorService(textRepo)
	textStatisticsQueryService := query.NewTextStatisticsQueryService(textRepo)
	textQueryService := query.NewTextQueryService(textRepo)
	ctx := context.Background()

	handler := transport.NewHandler(ctx, valuatorService, textStatisticsQueryService, textQueryService)

	r := mux.NewRouter()

	r.HandleFunc("/", handler.Index).Methods("GET")
	r.HandleFunc("/summary", handler.Summary).Methods("POST")
	r.HandleFunc("/about", handler.About).Methods("GET")

	log.Println("Starting server on :8082")
	http.ListenAndServe(":8082", r)
}
