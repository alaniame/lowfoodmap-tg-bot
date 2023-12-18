package main

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4"
	"github.com/joho/godotenv"
	"log"
	hand "lowfoodmap-tg-bot/internal/handler"
	repo "lowfoodmap-tg-bot/internal/repository"
	serv "lowfoodmap-tg-bot/internal/service"
	"net/http"
	"os"
)

func initHandler(db *pgx.Conn, handler *hand.Handler) http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/product",
		func(w http.ResponseWriter, r *http.Request) {
			// GetProduct (db, w, r)
		}).Methods(http.MethodGet)
	r.HandleFunc("/product",
		func(w http.ResponseWriter, r *http.Request) {
			// GetProductsCategory (db, w, r)
		}).Methods(http.MethodGet)
	return r
}

func main() {
	// config
	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading .env file: %s", err)
	}

	// connect to db
	conn, err := pgx.Connect(context.Background(), os.Getenv("DB_HOST"))
	if err != nil {
		log.Fatalf("cannot connect to database: %s", err)
	}
	defer conn.Close(context.Background())

	repository := repo.NewRepository(conn)
	service := serv.NewService(repository)
	handler := hand.NewHandler(service)

	// handle requests
	http.Handle("/", initHandler(conn, handler))
	contactHttpErr := http.ListenAndServe(":8080", nil)
	if contactHttpErr != nil {
		log.Fatalf("server startup error: %v\n", contactHttpErr)
	}

}
