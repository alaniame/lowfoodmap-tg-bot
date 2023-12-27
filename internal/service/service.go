package service

import (
	entity "github.com/alaniame/lowfoodmap-tg-bot"
	"github.com/alaniame/lowfoodmap-tg-bot/internal/repository"
	"mime/multipart"
)

type Product interface {
	GetProduct(name string) (*entity.Product, error)
	UploadData(file multipart.File) error
	stringToProduct(record []string) (entity.Product, error)
}

type Service struct {
	Product
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Product: NewProductService(repos.Product, repos.ProductCategory, repos.CarbType),
	}
}
