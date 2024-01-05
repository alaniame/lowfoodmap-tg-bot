package main

import (
	hand "github.com/alaniame/lowfoodmap-tg-bot/internal/handler"
	repo "github.com/alaniame/lowfoodmap-tg-bot/internal/repository"
	serv "github.com/alaniame/lowfoodmap-tg-bot/internal/service"
	"github.com/alaniame/lowfoodmap-tg-bot/pkg/postgres"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	pool, err := postgres.NewPostgresDB()
	if err != nil {
		log.Fatalf("failed to initialize db: %s", err.Error())
	}
	defer pool.Close()

	productRepository := repo.NewProductRepository(pool)
	productCategoryRepository := repo.NewProductCategoryRepository(pool)
	CarbTypeRepository := repo.NewCarbTypeRepository(pool)
	service := serv.NewProductService(productRepository, productCategoryRepository, CarbTypeRepository)
	handler := hand.NewHandler(service)

	bot, err := tgbotapi.NewBotAPI(os.Getenv("TG_API_KEY"))
	if err != nil {
		log.Fatalf("cannot start bot: %s", err)
	}
	updateConfig := tgbotapi.NewUpdate(0)
	updates := bot.GetUpdatesChan(updateConfig)
	for update := range updates {
		if update.Message != nil {
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			var ans string
			if update.Message.Text == "/start" {
				ans = "Добро пожаловать! Введите название продукта"
			} else {
				ans = handler.GetProduct(update.Message.Text)
			}
			ans = EscapeText(ans)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, ans)
			msg.ParseMode = tgbotapi.ModeMarkdownV2
			bot.Send(msg)
		}
	}

	http.Handle("/", handler.InitRoutes())
	contactHttpErr := http.ListenAndServe(":8080", nil)
	if contactHttpErr != nil {
		log.Fatalf("server startup error: %v\n", contactHttpErr)
	}
}

func EscapeText(text string) string {
	var replacer *strings.Replacer
	replacer = strings.NewReplacer("(", "\\(", ")", "\\)", "-", "\\-", ".", "\\.")
	return replacer.Replace(text)
}
