package repository

import (
	"context"
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
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS product_categories (
		category_id SERIAL PRIMARY KEY,
		category_name VARCHAR(255) NOT NULL UNIQUE
	);

	CREATE TABLE IF NOT EXISTS carb_types (
		carb_id SERIAL PRIMARY KEY,
		carb_name VARCHAR(255) NOT NULL UNIQUE
	);

	CREATE TABLE IF NOT EXISTS products (
		id SERIAL PRIMARY KEY,
		product_name VARCHAR(255) NOT NULL UNIQUE,
		category_id INT NOT NULL,
		stage INT NOT NULL,
		portion_high INT,
		portion_medium INT,
		portion_low INT,
		portion_size VARCHAR(255),
		FOREIGN KEY (category_id) REFERENCES product_categories (category_id)
	);
	
	CREATE TABLE IF NOT EXISTS product_carb_types (
	    id SERIAL PRIMARY KEY,
	  	product_id INT NOT NULL,
	  	carb_id INT NOT NULL,
	  	FOREIGN KEY (carb_id) REFERENCES carb_types(carb_id),
	  	FOREIGN KEY (product_id) REFERENCES products(id)
	);`
	_, err := conn.Exec(context.Background(), createTableSQL)
	if err != nil {
		log.Fatalf("error creating table: %v\n", err)
	}

	for carbName := range carbTypeMap {
		addInitialData := fmt.Sprintf("INSERT INTO carb_types (carb_name) VALUES ('%s') ON CONFLICT (carb_name) DO NOTHING;", carbName)
		_, err = conn.Exec(context.Background(), addInitialData)
		if err != nil {
			log.Fatalf("error adding carbTypes to table: %v\n", err)
		}
	}

	for categoryName := range productCategoryMap {
		addInitialData := fmt.Sprintf("INSERT INTO product_categories (category_name) VALUES ('%s') ON CONFLICT (category_name) DO NOTHING;", categoryName)
		_, err = conn.Exec(context.Background(), addInitialData)
		if err != nil {
			log.Fatalf("error adding productCategories to table: %v\n", err)
		}
	}

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
		addProduct := fmt.Sprintf(`
			INSERT INTO products (product_name, category_id, stage, portion_high, portion_medium, portion_low, portion_size)
			VALUES ('%s', '%d', '%d', '%d', '%d', '%d', '%s')  ON CONFLICT (product_name) DO NOTHING;`,
			product.ProductName,
			product.CategoryId,
			product.Stage,
			product.PortionHigh,
			product.PortionMedium,
			product.PortionLow,
			product.PortionSize)
		_, err := r.conn.Exec(context.Background(), addProduct)
		if err != nil {
			log.Fatalf("error adding product to table: %v\n", err)
		}

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
