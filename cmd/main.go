package main

import (
	"context"
	hand "github.com/alaniame/lowfoodmap-tg-bot/internal/handler"
	repo "github.com/alaniame/lowfoodmap-tg-bot/internal/repository"
	serv "github.com/alaniame/lowfoodmap-tg-bot/internal/service"
	"github.com/alaniame/lowfoodmap-tg-bot/pkg/postgres"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"net/http"
	"os"
)

func main() {
	pool, err := postgres.NewPostgresDB()
	if err != nil {
		log.Fatalf("failed to initialize db: %s", err.Error())
	}
	defer pool.Close()

	if err := pool.Ping(context.Background()); err != nil {
		log.Fatalf("cannot ping db: %s", err)
	}

	productRepository := repo.NewProductRepository(pool)
	productCategoryRepository := repo.NewProductCategoryRepository(pool)
	carbTypeRepository := repo.NewCarbTypeRepository(pool)
	service := serv.NewProductService(productRepository, productCategoryRepository, carbTypeRepository)
	handler := hand.NewHandler(service)

	bot, err := tgbotapi.NewBotAPI(os.Getenv("TG_API_KEY"))
	if err != nil {
		log.Fatalf("cannot start bot: %s", err)
	}
	updateConfig := tgbotapi.NewUpdate(0)
	updates := bot.GetUpdatesChan(updateConfig)

	done := make(chan struct{})

	go func() {
		for update := range updates {
			if update.Message == nil {
				continue
			}
			msg := hand.HandleMessage(update.Message, *handler)
			if _, err := bot.Send(msg); err != nil {
				log.Printf("send meessage eroor: %s", err)
			}
		}
		done <- struct{}{}
	}()

	go func() {
		http.Handle("/", handler.InitRoutes())
		contactHttpErr := http.ListenAndServe(":8080", nil)
		if contactHttpErr != nil {
			log.Fatalf("server startup error: %v\n", contactHttpErr)
		}
		done <- struct{}{}
	}()

	<-done
	<-done
}
