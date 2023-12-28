package main

import (
	hand "github.com/alaniame/lowfoodmap-tg-bot/internal/handler"
	repo "github.com/alaniame/lowfoodmap-tg-bot/internal/repository"
	serv "github.com/alaniame/lowfoodmap-tg-bot/internal/service"
	"github.com/alaniame/lowfoodmap-tg-bot/pkg/postgres"
	"log"
	"net/http"
)

func main() {
	conn, err := postgres.NewPostgresDB()
	if err != nil {
		log.Fatalf("failed to initialize db: %s", err.Error())
	}

	productRepository := repo.NewProductRepository(conn)
	productCategoryRepository := repo.NewProductCategoryRepository(conn)
	CarbTypeRepository := repo.NewCarbTypeRepository(conn)
	service := serv.NewProductService(productRepository, productCategoryRepository, CarbTypeRepository)
	handler := hand.NewHandler(service)

	http.Handle("/", handler.InitRoutes())
	contactHttpErr := http.ListenAndServe(":8080", nil)
	if contactHttpErr != nil {
		log.Fatalf("server startup error: %v\n", contactHttpErr)
	}

}
