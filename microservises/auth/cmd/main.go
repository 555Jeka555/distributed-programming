package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"server/pkg/app/service"
	"server/pkg/infrastructure/redis/hash"
	"server/pkg/infrastructure/redis/repo"
	"server/pkg/infrastructure/transport"
	"time"
)

func main() {
	mainRdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_AUTH_URL"),
		Password: os.Getenv("REDIS_AUTH_PASSWORD"),
		Username: "default",
	})

	pingCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := mainRdb.Ping(pingCtx).Result(); err != nil {
		log.Fatalf("Failed to connect to main Redis: %v", err)
	}
	log.Println("Successfully connected to main Redis")

	ctx := context.Background()

	userRepo := repo.NewUserRepository(mainRdb)
	hashService := hash.NewHashService()
	userService := service.NewUserService(hashService, userRepo)

	handler := transport.NewHandler(ctx, os.Getenv("JWT_KEY"), &userService)

	r := mux.NewRouter()

	r.HandleFunc("/auth/login-page", handler.LoginPage).Methods("GET")
	r.HandleFunc("/auth/login", handler.Login).Methods("POST")
	r.HandleFunc("/auth/register", handler.Register).Methods("POST")
	r.HandleFunc("/auth/logout", handler.Logout).Methods("POST")

	log.Println(fmt.Sprintf("Starting server on %s", os.Getenv("LISTENING_SERVER_PORT")))
	http.ListenAndServe(os.Getenv("LISTENING_SERVER_PORT"), r)
}
