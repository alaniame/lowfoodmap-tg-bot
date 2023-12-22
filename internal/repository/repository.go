package repository

import (
	"context"
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
	createTableSQL := `CREATE TABLE IF NOT EXISTS product_categories (
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
	_, createErr := conn.Exec(context.Background(), createTableSQL)
	if createErr != nil {
		log.Fatalf("error creating table: %v\n", createErr)
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
