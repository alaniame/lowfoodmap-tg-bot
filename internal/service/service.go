package service

import (
	"encoding/csv"
	"github.com/alaniame/lowfoodmap-tg-bot/internal/repository"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"strconv"
)

type Service struct {
	repo repository.Repository
}

func NewService(repo repository.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetProduct(name string) (*repository.Product, error) {
	return s.repo.GetProduct(name)
}

func (s *Service) UploadData(w http.ResponseWriter, file multipart.File) {
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
	s.repo.AddProducts(products)
}
