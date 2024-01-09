package handler

import (
	"github.com/alaniame/lowfoodmap-tg-bot/internal/service"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strings"
)

type Handler struct {
	service *service.ProductService
}

func NewHandler(service *service.ProductService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) InitRoutes() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/upload_data",
		h.UploadData).Methods(http.MethodPost)
	return r
}

func HandleMessage(message *tgbotapi.Message, handler Handler) tgbotapi.MessageConfig {
	log.Printf("[%s] %s", message.From.UserName, message.Text)
	var ans string
	if message.Text == "/start" {
		ans = "Добро пожаловать! Введите название продукта"
	} else {
		ans = handler.GetProduct(message.Text)
	}
	ans = EscapeText(ans)
	msg := tgbotapi.NewMessage(message.Chat.ID, ans)
	msg.ParseMode = tgbotapi.ModeMarkdownV2
	return msg
}

func EscapeText(text string) string {
	replacer := strings.NewReplacer("(", "\\(", ")", "\\)", "-", "\\-", ".", "\\.")
	return replacer.Replace(text)
}
