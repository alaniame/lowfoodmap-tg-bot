package handler

import (
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
)

func (h *Handler) GetProduct(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "product name is empty", http.StatusBadRequest)
		return
	}
	product, err := h.service.GetProduct(name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
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
			log.Printf("file close error: %v\n", err)
		}
	}(file)
	err = h.service.UploadData(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
