package main

import (
	"lowfoodmap-tg-bot/internal/handler"
	"lowfoodmap-tg-bot/internal/repository"
	"lowfoodmap-tg-bot/internal/service"
)

func main() {
	repo := repository.NewRepository
	service := service.NewService(repo)
	handler := handler.NewHandler(service)
}
