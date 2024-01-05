package handler

import (
	"github.com/alaniame/lowfoodmap-tg-bot/internal/service"
	"github.com/gorilla/mux"
	"net/http"
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
