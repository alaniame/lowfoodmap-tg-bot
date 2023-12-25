package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4"
	"log"
)

type Repository struct {
	conn *pgx.Conn
}

type Product struct {
	ProductName   string
	PortionHigh   int
	PortionMedium int
	PortionLow    int
	PortionSize   string
	CarbId        []CarbType
	Stage         int
	CategoryId    ProductCategory
}

func NewRepository(conn *pgx.Conn) Repository {
	return Repository{conn: conn}
}

func (r *Repository) GetProduct(name string) (*Product, error) {
	row := r.conn.QueryRow(context.Background(),
		"SELECT product_name FROM products WHERE product_name = $1;",
		name)
	var prod Product
	err := row.Scan(&prod.ProductName)
	return &prod, err
}

func (r *Repository) AddProducts(products []Product) {
	for _, product := range products {
		addProduct := `INSERT INTO products (product_name, category_id, stage, portion_high, portion_medium, portion_low, portion_size)
			VALUES ($1, $2, $3, $4, $5, $6, $7)  ON CONFLICT (product_name) DO NOTHING RETURNING id;`
		var productId int
		err := r.conn.QueryRow(context.Background(), addProduct, product.ProductName, product.CategoryId, product.Stage, product.PortionHigh, product.PortionMedium, product.PortionLow, product.PortionSize).Scan(&productId)
		if err != nil {
			if !errors.Is(err, pgx.ErrNoRows) {
				log.Fatalf("error adding product to table: %v\n", err)
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
				_, err = r.conn.Exec(context.Background(), addCarbTypeRelation)
				if err != nil {
					log.Fatalf("error adding carb type relation to table: %v\n", err)
				}
			}
		}
	}
}
