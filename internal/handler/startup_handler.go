package handler

import(
	"encoding/json"
	"net/http"
	"github.com/rudraa2005/mic-website-main/backend/internal/model"
    "github.com/rudraa2005/mic-website-main/backend/internal/service"
)

type StartupHandler struct {
    svc *service.StartupService
}

func NewStartupHandler(s *service.StartupService) *StartupHandler {
    return &StartupHandler{s}
}

func (h *StartupHandler) Create(w http.ResponseWriter, r *http.Request) {
    var input model.Startup
    if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
        http.Error(w, "invalid JSON", http.StatusBadRequest)
        return 
    }

    input.OwnerID = "dummy-user"

    err := h.svc.CreateStartup(r.Context(), &input)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(input)
}

func (h *StartupHandler) GetMine(w http.ResponseWriter, r *http.Request) {
    ownerID := "dummy-user"
    startups, err := h.svc.ListMine(r.Context(), ownerID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    json.NewEncoder(w).Encode(startups)
}
