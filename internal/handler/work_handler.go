package handler

import (
	"encoding/json"
	"net/http"

	"github.com/rudraa2005/mic-website-main/backend/internal/repository"
)

type WorkHandler struct {
	repo *repository.SubmissionsRepo
}

func NewWorkHandler(repo *repository.SubmissionsRepo) *WorkHandler {
	return &WorkHandler{repo: repo}
}

func (h *WorkHandler) GetIncubationPipeline(w http.ResponseWriter, r *http.Request) {
	items, err := h.repo.GetIncubationPipeline(r.Context())
	if err != nil {
		http.Error(w, "failed to fetch incubation pipeline", 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if items == nil {
		w.Write([]byte("[]"))
		return
	}
	json.NewEncoder(w).Encode(items)
}
