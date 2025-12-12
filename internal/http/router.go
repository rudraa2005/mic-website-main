package http

import(
	"net/http"
	"github.com/go-chi/chi/v5"
	"github.com/rudraa2005/mic-website-main/backend/internal/handler"
)

func NewRouter(
    sh *handler.StartupHandler,
    rh *handler.ReviewHandler,
) *chi.Mux {

    r := chi.NewRouter()

    r.Post("/api/startups", sh.Create)
    r.Get("/api/startups/mine", sh.GetMine)

    r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("OK"))
    })

    r.Route("/api/reviews", func(rt chi.Router) {
        rt.Get("/mine", rh.ListMyReviews)
        rt.Post("/", rh.SubmitReview)
    })

    r.Get("/api/startups/{id}/reviews", rh.ListForStartup)

    return r
}

