package repository

import (
	"github.com/alaniame/lowfoodmap-tg-bot/internal/entity"
	"github.com/jackc/pgx/v4/pgxpool"
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

func NewCarbTypeRepository(conn *pgxpool.Pool) *CarbTypeRepository {
	return &CarbTypeRepository{pool: conn}
}

func NewProductCategoryRepository(conn *pgxpool.Pool) *ProductCategoryRepository {
	return &ProductCategoryRepository{pool: conn}
}

func NewProductRepository(conn *pgxpool.Pool) *ProductRepository {
	return &ProductRepository{pool: conn}
}
