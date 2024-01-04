package repository

import (
	"github.com/alaniame/lowfoodmap-tg-bot/internal/entity"
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
	GetProduct(productName string) ([]entity.ProductOutput, error)
}

func NewCarbTypeRepository(conn *pgx.Conn) *CarbTypeRepository {
	return &CarbTypeRepository{conn: conn}
}

func NewProductCategoryRepository(conn *pgx.Conn) *ProductCategoryRepository {
	return &ProductCategoryRepository{conn: conn}
}

func NewProductRepository(conn *pgx.Conn) *ProductRepository {
	return &ProductRepository{conn: conn}
}
