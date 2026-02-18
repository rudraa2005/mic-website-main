package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rudraa2005/mic-website-main/backend/internal/service"
)

type AdminFacultyHandler struct {
	service *service.AdminFacultyService
}

func NewAdminFacultyHandler(s *service.AdminFacultyService) *AdminFacultyHandler {
	return &AdminFacultyHandler{service: s}
}

// GetAllFaculty returns all users with FACULTY role
func (h *AdminFacultyHandler) GetAllFaculty(w http.ResponseWriter, r *http.Request) {
	faculty, err := h.service.GetAllFaculty(r.Context())
	if err != nil {
		http.Error(w, "failed to fetch faculty", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(faculty)
}

// CreateFaculty creates a new faculty user
func (h *AdminFacultyHandler) CreateFaculty(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" || req.Email == "" || req.Password == "" {
		http.Error(w, "name, email and password are required", http.StatusBadRequest)
		return
	}

	err := h.service.CreateFaculty(r.Context(), req.Name, req.Email, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// UpdateFaculty updates an existing faculty member
func (h *AdminFacultyHandler) UpdateFaculty(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "faculty id required", http.StatusBadRequest)
		return
	}

	var req struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	err := h.service.UpdateFaculty(r.Context(), id, req.Name, req.Email, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DeleteFaculty removes a faculty member
func (h *AdminFacultyHandler) DeleteFaculty(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "faculty id required", http.StatusBadRequest)
		return
	}

	err := h.service.DeleteFaculty(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
