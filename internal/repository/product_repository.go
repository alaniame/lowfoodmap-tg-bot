package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/alaniame/lowfoodmap-tg-bot/internal/entity"
	"github.com/jackc/pgx/v4"
	"log"
)

type ProductRepository struct {
	conn *pgx.Conn
}

func (r ProductRepository) AddProducts(products []entity.Product) error {
	for _, product := range products {
		tx, err := r.conn.Begin(context.Background())
		if err != nil {
			return fmt.Errorf("error adding transaction: %v\n", err)
		}
		addProduct := `INSERT INTO products (product_name, category_id, stage, portion_high, portion_medium, portion_low, portion_size)
			VALUES ($1, $2, $3, $4, $5, $6, $7)  ON CONFLICT (product_name) DO NOTHING RETURNING id;`
		var productId int
		err = tx.QueryRow(context.Background(), addProduct, product.ProductName, product.CategoryId, product.Stage, product.PortionHigh, product.PortionMedium, product.PortionLow, product.PortionSize).Scan(&productId)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				// Логируем, что продукт уже существует, но не прерываем выполнение
				log.Printf("product already exists: %s\n", product.ProductName)
			} else {
				err := tx.Rollback(context.Background())
				if err != nil {
					return fmt.Errorf("rollback error: %v\n", err)
				}
				return fmt.Errorf("error adding product to table: %v\n", err)
			}
		}

		if productId != 0 {
			for _, carbId := range product.CarbId {
				addCarbTypeRelation := fmt.Sprintf(`INSERT INTO product_carb_types 
				(product_id, carb_id) VALUES (currval('products_id_seq'), '%d')
				ON CONFLICT DO NOTHING;`, carbId)
				_, err = tx.Exec(context.Background(), addCarbTypeRelation)
				if err != nil {
					err := tx.Rollback(context.Background())
					if err != nil {
						return fmt.Errorf("rollback error: %v\n", err)
					}
					return fmt.Errorf("error adding carb type relation to table: %v\n", err)
				}
			}
		}
		err = tx.Commit(context.Background())
		if err != nil {
			return err
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
	rows, err := r.conn.Query(context.Background(), query, searchPattern)
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
