package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rudraa2005/mic-website-main/backend/internal/middleware"
	"github.com/rudraa2005/mic-website-main/backend/internal/service"
)

type FacultyProgressHandler struct {
	service *service.FacultyProgressService
}

func NewFacultyProgressHandler(
	service *service.FacultyProgressService,
) *FacultyProgressHandler {
	return &FacultyProgressHandler{service: service}
}

func (h *FacultyProgressHandler) GetMyProgress(
	w http.ResponseWriter,
	r *http.Request,
) {
	log.Println("HANDLER HIT")
	claims, err := middleware.GetUser(r)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	if claims.Role != "FACULTY" {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	data, err := h.service.GetProgressForFaculty(
		r.Context(),
		claims.UserID,
	)
	log.Println(data)
	if err != nil {
		http.Error(w, "failed to fetch progress", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func (h *FacultyProgressHandler) GetProgressBySubmission(
	w http.ResponseWriter,
	r *http.Request,
) {
	claims, err := middleware.GetUser(r)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	if claims.Role != "FACULTY" {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	submissionID := chi.URLParam(r, "submission_id")
	if submissionID == "" {
		http.Error(w, "missing submission id", http.StatusBadRequest)
		return
	}

	data, err := h.service.GetProgressBySubmission(
		r.Context(),
		claims.UserID,
		submissionID,
	)
	if err != nil {
		http.Error(w, "idea not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
