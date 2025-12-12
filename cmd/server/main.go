package main

import (
	"log"
	"net/http"
    "os"
	"github.com/rudraa2005/mic-website-main/backend/internal/db"
	"github.com/rudraa2005/mic-website-main/backend/internal/repository"
	"github.com/rudraa2005/mic-website-main/backend/internal/service"
	h "github.com/rudraa2005/mic-website-main/backend/internal/handler"
	r "github.com/rudraa2005/mic-website-main/backend/internal/http"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")

pool, err := db.NewPool(dsn)
if err != nil {
    log.Fatal("DB connection failed: ", err)
}
defer pool.Close()


	startupRepo := repository.NewStartupRepository(pool)
	reviewRepo := repository.NewReviewRepository(pool)

	startupService := service.NewStartupService(startupRepo)
	reviewService := service.NewReviewService(reviewRepo)

	startupHandler := h.NewStartupHandler(startupService)
	reviewHandler := h.NewReviewHandler(reviewService)

	router := r.NewRouter(startupHandler, reviewHandler)

	log.Println("Server running on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}
