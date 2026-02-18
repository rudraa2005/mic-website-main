package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rudraa2005/mic-website-main/backend/internal/middleware"
	"github.com/rudraa2005/mic-website-main/backend/internal/service"
)

type FacultyReviewHandler struct {
	service *service.FacultyReviewService
}

func NewFacultyReviewHandler(s *service.FacultyReviewService) *FacultyReviewHandler {
	return &FacultyReviewHandler{service: s}
}

func (h *FacultyReviewHandler) GetSubmitted(w http.ResponseWriter, r *http.Request) {
	items, err := h.service.GetSubmitted(r.Context())
	if err != nil {
		http.Error(w, "failed to fetch submissions", 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if items == nil {
		w.Write([]byte("[]"))
		return
	}
	json.NewEncoder(w).Encode(items)
}

func (h *FacultyReviewHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "missing id", 400)
		return
	}

	item, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, "failed to fetch submission", 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}

func (h *FacultyReviewHandler) Decide(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	claims, err := middleware.GetUser(r)

	if err != nil || claims.Role != "FACULTY" {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	var body struct {
		Decision string `json:"decision"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid request", 400)
		return
	}

	err = h.service.Decide(r.Context(), id, claims.UserID, body.Decision)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
