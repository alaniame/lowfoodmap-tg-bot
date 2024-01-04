package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/alaniame/lowfoodmap-tg-bot/internal/entity"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
)

type ProductRepository struct {
	pool *pgxpool.Pool
}

func (r ProductRepository) AddProducts(products []entity.Product) error {
	batch := &pgx.Batch{}
	for _, product := range products {
		batch.Queue(`INSERT INTO products (product_name, category_id, stage, portion_high, portion_medium, portion_low, portion_size)
                VALUES ($1, $2, $3, $4, $5, $6, $7) ON CONFLICT (product_name) DO NOTHING RETURNING id;`,
			product.ProductName, product.CategoryId, product.Stage, product.PortionHigh, product.PortionMedium, product.PortionLow, product.PortionSize)
	}
	res := r.pool.SendBatch(context.Background(), batch)
	defer func(res pgx.BatchResults) {
		err := res.Close()
		if err != nil {
			log.Printf("error closing pool: %s\n", err)
		}
	}(res)
	var productIds []int
	for i := 0; i < batch.Len(); i++ {
		var productId int
		err := res.QueryRow().Scan(&productId)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				continue
			}
			return fmt.Errorf("error adding product to table: %v\n", err)
		}
		productIds = append(productIds, productId)
	}
	for i, product := range products {
		for _, carbId := range product.CarbId {
			if i < len(productIds) {
				batch.Queue(`INSERT INTO product_carb_types (product_id, carb_id) VALUES ($1, $2) ON CONFLICT DO NOTHING;`,
					productIds[i], carbId)
			}
		}
	}
	res = r.pool.SendBatch(context.Background(), batch)
	defer func(res pgx.BatchResults) {
		err := res.Close()
		if err != nil {
			log.Printf("error closing pool: %s\n", err)
		}
	}(res)
	for i := 0; i < batch.Len(); i++ {
		_, err := res.Exec()
		if err != nil {
			return fmt.Errorf("error adding carb type relation to table: %v\n", err)
		}
	}
	return nil
}

func (r ProductRepository) GetProduct(productName string) ([]entity.ProductOutput, error) {
	query := `SELECT
        p.product_name,
        p.stage,
        p.portion_high,
        p.portion_medium,
        p.portion_low,
        p.portion_size,
        string_agg(COALESCE(ct.carb_name, ''), ', ') AS carb_names
    FROM products p
        LEFT JOIN product_carb_types pct ON p.id = pct.product_id
        LEFT JOIN carb_types ct ON pct.carb_id = ct.carb_id
    WHERE LOWER(p.product_name) LIKE LOWER($1)
    GROUP BY p.product_name, p.stage, p.portion_high, p.portion_medium, p.portion_low, p.portion_size
    ORDER BY p.stage, p.product_name;`
	searchPattern := "%" + productName + "%"
	rows, err := r.pool.Query(context.Background(), query, searchPattern)
	var prodOuts []entity.ProductOutput
	if err != nil {
		return prodOuts, err
	}
	defer rows.Close()
	for rows.Next() {
		var prodOut entity.ProductOutput
		err := rows.Scan(&prodOut.ProductName, &prodOut.Stage, &prodOut.PortionHigh, &prodOut.PortionMedium, &prodOut.PortionLow, &prodOut.PortionSize, &prodOut.Carbs)
		if err != nil {
			return prodOuts, err
		}
		prodOuts = append(prodOuts, prodOut)
	}
	if err = rows.Err(); err != nil {
		return prodOuts, err
	}

	return prodOuts, nil
}
