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
	ProductName string
}

func NewRepository(conn *pgx.Conn) Repository {
	createTableSQL := `CREATE TABLE IF NOT EXISTS products (
		id SERIAL PRIMARY KEY,
		product_name VARCHAR(50),
		stage INT
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
