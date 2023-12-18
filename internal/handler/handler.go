package handler

import (
	"lowfoodmap-tg-bot/internal/service"
	"net/http"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetProduct(w http.ResponseWriter, r *http.Request) {
}
