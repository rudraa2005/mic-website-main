package main

import(
	"log"
	"net/http"
	"github.com/rudraa2005/mic-website-main/backend/internal/db"
    "github.com/rudraa2005/mic-website-main/backend/internal/repository"
    "github.com/rudraa2005/mic-website-main/backend/internal/service"
    h "github.com/rudraa2005/mic-website-main/backend/internal/handler"
    r "github.com/rudraa2005/mic-website-main/backend/internal/http"
)

func main() {
    pool, err := db.NewPool()
    if err != nil {
        log.Fatal("DB connection failed: ", err)
    }

    startupRepo := repository.NewStartupRepository(pool)
    startupService := service.NewStartupService(startupRepo)
    startupHandler := h.NewStartupHandler(startupService)

    router := r.NewRouter(startupHandler)

    log.Println("Server running on :8080")
    http.ListenAndServe(":8080", router)


}