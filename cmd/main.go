package main

import (
	hand "github.com/alaniame/lowfoodmap-tg-bot/internal/handler"
	"github.com/alaniame/lowfoodmap-tg-bot/internal/postgres"
	repo "github.com/alaniame/lowfoodmap-tg-bot/internal/repository"
	serv "github.com/alaniame/lowfoodmap-tg-bot/internal/service"
	"log"
	"net/http"
)

func main() {
	conn, err := postgres.NewPostgresDB()
	if err != nil {
		log.Fatalf("failed to initialize db: %s", err.Error())
	}

	repository := repo.NewRepository(conn)
	service := serv.NewService(repository)
	handler := hand.NewHandler(service)

	http.Handle("/", handler.InitRoutes())
	contactHttpErr := http.ListenAndServe(":8080", nil)
	if contactHttpErr != nil {
		log.Fatalf("server startup error: %v\n", contactHttpErr)
	}

}
