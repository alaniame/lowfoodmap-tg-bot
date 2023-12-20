package main

import (
	"context"
	"fmt"
	hand "github.com/alaniame/lowfoodmap-tg-bot/internal/handler"
	repo "github.com/alaniame/lowfoodmap-tg-bot/internal/repository"
	serv "github.com/alaniame/lowfoodmap-tg-bot/internal/service"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"strings"
)

func initHandler(handler *hand.Handler) http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/product",
		handler.GetProduct).Methods(http.MethodGet)
	return r
}

func example() {
	unusedVar := 5
}

func main() {
	// config
	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading .env file: %s", err)
	}
	dbLogin := os.Getenv("POSTGRES_USER")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")
	hostPort := strings.Split(os.Getenv("POSTGRES_PORT"), ":")[0]
	dbURL := fmt.Sprintf("postgres://%s:%s@localhost:%s/%s", dbLogin, dbPassword, hostPort, dbName)

	// connect to db
	conn, err := pgx.Connect(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("cannot connect to database: %s", err)
	}
	defer func(conn *pgx.Conn, ctx context.Context) {
		err := conn.Close(ctx)
		if err != nil {
			log.Fatalf("db close error: %v\n", err)
		}
	}(conn, context.Background())

	// layers
	repository := repo.NewRepository(conn)
	service := serv.NewService(repository)
	handler := hand.NewHandler(service)

	// handle requests
	http.Handle("/", initHandler(handler))
	contactHttpErr := http.ListenAndServe(":8080", nil)
	if contactHttpErr != nil {
		log.Fatalf("server startup error: %v\n", contactHttpErr)
	}

}
