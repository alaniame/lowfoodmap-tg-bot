package handler

import (
	"encoding/csv"
	"fmt"
	"github.com/alaniame/lowfoodmap-tg-bot/internal/repository"
	"github.com/alaniame/lowfoodmap-tg-bot/internal/service"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"strconv"
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
	_, err = w.Write([]byte(responseString))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) UploadData(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(1 << 20); err != nil { // 1 MB
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {
			log.Fatalf("file close error: %v\n", err)
		}
	}(file)
	csvReader := csv.NewReader(file)
	var products []repository.Product
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		portionHigh, err := strconv.Atoi(record[1])
		if err != nil {
			portionHigh = 0
		}
		portionMedium, err := strconv.Atoi(record[2])
		if err != nil {
			portionHigh = 0
		}
		portionLow, err := strconv.Atoi(record[3])
		if err != nil {
			portionHigh = 0
		}
		carbTypes, err := repository.StringToCarbTypes(record[5])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Fatalf("error converting CarbTypes: %v", err)
			return
		}
		stage, err := strconv.Atoi(record[6])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Fatalf("error converting Stage: %v", err)
			return
		}
		category, err := repository.StringToProductCategory(record[7])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Fatalf("error converting Category: %v", err)
			return
		}
		product := repository.Product{
			ProductName:   record[0],
			PortionHigh:   portionHigh,
			PortionMedium: portionMedium,
			PortionLow:    portionLow,
			PortionSize:   record[4],
			CarbId:        carbTypes,
			Stage:         stage,
			CategoryId:    category,
		}
		products = append(products, product)
	}
	h.service.UploadData(products)
}
