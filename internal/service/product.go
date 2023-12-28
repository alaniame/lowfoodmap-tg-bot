package service

import (
	"encoding/csv"
	"fmt"
	"github.com/alaniame/lowfoodmap-tg-bot/internal/entity"
	"github.com/alaniame/lowfoodmap-tg-bot/internal/repository"
	"io"
	"mime/multipart"
	"strconv"
)

type ProductService struct {
	product         repository.Product
	productCategory repository.ProductCategory
	carbType        repository.CarbType
}

func (s *ProductService) GetProduct(name string) ([]entity.ProductOutput, error) {
	return s.product.GetProduct(name)
}

func (s *ProductService) UploadData(file multipart.File) error {
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
		product, err := s.stringToProduct(record)
		if err != nil {
			return err
		}
		products = append(products, product)
	}
	err := s.product.AddProducts(products)
	if err != nil {
		return err
	}
	return nil
}

func (s *ProductService) stringToProduct(record []string) (entity.Product, error) {
	var product entity.Product
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
	carbTypes, err := s.carbType.GetCarbIds(record[5])
	if err != nil {
		return product, fmt.Errorf("error converting CarbTypes: %v\n", err)
	}
	stage, err := strconv.Atoi(record[6])
	if err != nil {
		return product, fmt.Errorf("error converting Stage: %v\n", err)
	}
	category, err := s.productCategory.GetProductCategoryId(record[7])
	if err != nil {
		return product, fmt.Errorf("error converting Category: %v\n", err)
	}
	product = entity.Product{
		ProductName:   record[0],
		PortionHigh:   portionHigh,
		PortionMedium: portionMedium,
		PortionLow:    portionLow,
		PortionSize:   record[4],
		CarbId:        carbTypes,
		Stage:         stage,
		CategoryId:    category,
	}
	return product, nil
}
