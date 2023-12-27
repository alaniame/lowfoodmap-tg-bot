package service

import (
	"encoding/csv"
	"fmt"
	"github.com/alaniame/lowfoodmap-tg-bot"
	"github.com/alaniame/lowfoodmap-tg-bot/internal/repository"
	"io"
	"mime/multipart"
	"strconv"
)

type Service struct {
	repo *repository.Repository
}

func NewService(repo *repository.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetProduct(name string) (*entity.Product, error) {
	return s.repo.GetProduct(name)
}

func (s *Service) UploadData(file multipart.File) error {
	csvReader := csv.NewReader(file)
	var products []entity.Product
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
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
		carbTypes, err := s.repo.GetCarbIds(record[5])
		if err != nil {
			return fmt.Errorf("error converting CarbTypes: %v\n", err)
		}
		stage, err := strconv.Atoi(record[6])
		if err != nil {
			return fmt.Errorf("error converting Stage: %v\n", err)
		}
		category, err := s.repo.GetProductCategoryId(record[7])
		if err != nil {
			return fmt.Errorf("error converting Category: %v\n", err)
		}
		product := entity.Product{
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
	err := s.repo.AddProducts(products)
	if err != nil {
		return err
	}
	return nil
}
