package repository

import (
	"github.com/alaniame/lowfoodmap-tg-bot"
	"github.com/jackc/pgx/v4"
)

type CarbType interface {
	GetCarbIds(carbName string) ([]int, error)
}

type ProductCategory interface {
	GetProductCategoryId(categoryName string) (int, error)
}

type Product interface {
	AddProducts(products []entity.Product) error
	GetProduct(productName string) (*entity.Product, error)
}

type Repository struct {
	CarbType
	ProductCategory
	Product
}

func NewRepository(conn *pgx.Conn) *Repository {
	return &Repository{
		CarbType:        NewCarbTypePostgres(conn),
		ProductCategory: NewProductCategoryPostgres(conn),
		Product:         NewProductPostgres(conn),
	}
}
