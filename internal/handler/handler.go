package handler

import (
	"fmt"
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
	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "product name is empty", http.StatusBadRequest)
		return
	}
	product, err := h.service.GetProduct(name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	responseString := fmt.Sprintf("Продукт: %s", product.ProductName)
	w.Write([]byte(responseString))
}
