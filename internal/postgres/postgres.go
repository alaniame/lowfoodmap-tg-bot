package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/joho/godotenv"
	"os"
	"strings"
)

func NewPostgresDB() (*pgx.Conn, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}
	dbLogin := os.Getenv("POSTGRES_USER")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")
	hostPort := strings.Split(os.Getenv("POSTGRES_PORT"), ":")[0]
	dbURL := fmt.Sprintf("postgres://%s:%s@localhost:%s/%s", dbLogin, dbPassword, hostPort, dbName)

	conn, err := pgx.Connect(context.Background(), dbURL)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
