package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/alaniame/lowfoodmap-tg-bot"
	"github.com/jackc/pgx/v4"
	"log"
)

type ProductPostgres struct {
	conn *pgx.Conn
}

func NewProductPostgres(conn *pgx.Conn) *ProductPostgres {
	return &ProductPostgres{conn: conn}
}

func (r ProductPostgres) AddProducts(products []entity.Product) {
	for _, product := range products {
		tx, err := r.conn.Begin(context.Background())
		if err != nil {
			log.Fatalf("error adding transaction: %v\n", err)
			return
		}
		addProduct := `INSERT INTO products (product_name, category_id, stage, portion_high, portion_medium, portion_low, portion_size)
			VALUES ($1, $2, $3, $4, $5, $6, $7)  ON CONFLICT (product_name) DO NOTHING RETURNING id;`
		var productId int
		err = tx.QueryRow(context.Background(), addProduct, product.ProductName, product.CategoryId, product.Stage, product.PortionHigh, product.PortionMedium, product.PortionLow, product.PortionSize).Scan(&productId)
		if err != nil {
			if !errors.Is(err, pgx.ErrNoRows) {
				err := tx.Rollback(context.Background())
				if err != nil {
					log.Fatalf("rollback error: %v\n", err)
					return
				}
				log.Fatalf("error adding product to table: %v\n", err)
				return
			} else {
				// Логируем, что продукт уже существует, но не прерываем выполнение
				log.Printf("product already exists: %s\n", product.ProductName)
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
						log.Fatalf("rollback error: %v\n", err)
						return
					}
					log.Fatalf("error adding carb type relation to table: %v\n", err)
					return
				}
			}
		}
		err = tx.Commit(context.Background())
		if err != nil {
			return
		}
	}
}

func (r ProductPostgres) GetProduct(productName string) (*entity.Product, error) {
	query := `SELECT product_name FROM products WHERE product_name = $1;`
	row := r.conn.QueryRow(context.Background(), query, productName)
	var prod entity.Product
	err := row.Scan(&prod.ProductName)
	return &prod, err
}