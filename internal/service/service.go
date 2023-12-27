package service

import (
	"github.com/alaniame/lowfoodmap-tg-bot/internal/entity"
	"github.com/alaniame/lowfoodmap-tg-bot/internal/repository"
	"mime/multipart"
)

type Product interface {
	GetProduct(name string) (*entity.Product, error)
	UploadData(file multipart.File) error
	stringToProduct(record []string) (entity.Product, error)
}

func NewProductService(product repository.Product, productCategory repository.ProductCategory, carbType repository.CarbType) *ProductService {
	return &ProductService{product: product, productCategory: productCategory, carbType: carbType}
}
