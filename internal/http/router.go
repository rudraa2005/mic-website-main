package http

import(
	"net/http"
	"github.com/go-chi/chi/v5"
	"github.com/rudraa2005/mic-website-main/backend/internal/handler"
)

func NewRouter(sh *handler.StartupHandler) *chi.Mux {
    r := chi.NewRouter()

    r.Post("/api/startups", sh.Create)
    r.Get("/api/startups/mine", sh.GetMine)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
})

    return r
}

