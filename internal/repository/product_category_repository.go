package repository

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
)

type ProductCategoryRepository struct {
	pool *pgxpool.Pool
}

func (r ProductCategoryRepository) GetProductCategoryId(categoryName string) (int, error) {
	query := `SELECT category_id FROM product_categories WHERE category_name = $1;`
	row := r.pool.QueryRow(context.Background(), query, categoryName)
	var id int
	err := row.Scan(&id)
	return id, err
}
