package handler

import (
	"fmt"
	"github.com/alaniame/lowfoodmap-tg-bot/internal/entity"
	"log"
	"mime/multipart"
	"net/http"
	"strings"
	"unicode/utf8"
)

func (h *Handler) GetProduct(name string) string {
	defaultError := "Что-то пошло не так, попробуйте еще раз. Мы уже занимаемся изучением проблемы"
	name = strings.TrimSpace(name)
	if name == "" {
		log.Println("product name is empty")
		return defaultError
	}
	length := utf8.RuneCountInString(name)
	if length < 3 {
		log.Println("too short name")
		return "Введите хотя бы 3 буквы"
	}
	products, err := h.service.GetProduct(name)
	if err != nil {
		log.Println(err.Error())
		return defaultError
	}
	if len(products) == 0 {
		log.Println("no products")
		return "Продукт не найден"
	}
	var responseBuilder strings.Builder
	for _, product := range products {
		productString := formatProductResponse(product)
		responseBuilder.WriteString(productString)
	}
	res := responseBuilder.String()
	return res
}

func formatProductResponse(product entity.ProductOutput) string {
	var responseBuilder strings.Builder
	responseBuilder.WriteString(fmt.Sprintf("Название: *%s*\nЭтап: *%d*\n", product.ProductName, product.Stage))
	if product.PortionHigh != 0 {
		responseBuilder.WriteString(fmt.Sprintf("🔴 Порция с высоким содержанием веществ FODMAP: *%d грамм*\n", product.PortionHigh))
	}
	if product.PortionMedium != 0 {
		responseBuilder.WriteString(fmt.Sprintf("🟡 Порция с умеренным содержанием веществ FODMAP: *%d грамм*\n", product.PortionMedium))
	}
	if product.PortionLow != 0 {
		responseBuilder.WriteString(fmt.Sprintf("🟢 Порция с низким содержанием веществ FODMAP: *%d грамм*\n", product.PortionLow))
	}
	if product.PortionSize != "" {
		responseBuilder.WriteString(fmt.Sprintf("Средний размер разрешенной порции: *%s*\n", product.PortionSize))
	}
	if product.Carbs != "" {
		responseBuilder.WriteString(fmt.Sprintf("Углеводы: *%s*\n", product.Carbs))
	}
	responseBuilder.WriteString("\n\n")
	return responseBuilder.String()
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
			log.Printf("file close error: %v\n", err)
		}
	}(file)
	err = h.service.UploadData(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
