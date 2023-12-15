package repository

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Repository struct {
	db *pgx.Conn
}

func NewRepository(db *pgx.Conn) *Repository {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading .env file: %s", err)
	}
	conn, err := pgx.Connect(context.Background(), os.Getenv("DB_HOST"))
	if err != nil {
		log.Fatalf("cannot connect to database: %s", err)
	}
	defer conn.Close(context.Background())
	return &Repository{db: db}
}
