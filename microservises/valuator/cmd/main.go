package main

import (
	"context"
	"fmt"
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
	textQueryService := query.NewTextQueryService(textRepo)
	ctx := context.Background()

	handler := transport.NewHandler(ctx, valuatorService, textQueryService)

	r := mux.NewRouter()

	r.HandleFunc("/", handler.Index).Methods("GET")
	r.HandleFunc("/summary", handler.Summary).Methods("POST")
	r.HandleFunc("/about", handler.About).Methods("GET")

	log.Println(fmt.Sprintf("Starting server on %s", os.Getenv("LISTENING_SERVER_PORT")))
	http.ListenAndServe(os.Getenv("LISTENING_SERVER_PORT"), r)
}
