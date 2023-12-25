package main

import (
	hand "github.com/alaniame/lowfoodmap-tg-bot/internal/handler"
	repo "github.com/alaniame/lowfoodmap-tg-bot/internal/repository"
	serv "github.com/alaniame/lowfoodmap-tg-bot/internal/service"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func initHandler(handler *hand.Handler) http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/product",
		handler.GetProduct).Methods(http.MethodGet)
	r.HandleFunc("/upload_data",
		handler.UploadData).Methods(http.MethodPost)
	return r
}

func main() {
	conn, err := repo.NewPostgresDB()
	if err != nil {
		log.Fatalf("failed to initialize db: %s", err.Error())
	}

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
