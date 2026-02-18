package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rudraa2005/mic-website-main/backend/internal/middleware"
	"github.com/rudraa2005/mic-website-main/backend/internal/repository"
	"github.com/rudraa2005/mic-website-main/backend/internal/service"
)

type FacultyIncubationHandler struct {
	progressService *service.FacultyProgressService
	companyRepo     *repository.CompanyRepo
}

func NewFacultyIncubationHandler(ps *service.FacultyProgressService, cr *repository.CompanyRepo) *FacultyIncubationHandler {
	return &FacultyIncubationHandler{
		progressService: ps,
		companyRepo:     cr,
	}
}

func (h *FacultyIncubationHandler) GetPortfolio(w http.ResponseWriter, r *http.Request) {
	claims, err := middleware.GetUser(r)
	if err != nil {
		http.Error(w, "unauthorized", 401)
		return
	}

	data, err := h.progressService.GetProgressForFaculty(r.Context(), claims.UserID)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func (h *FacultyIncubationHandler) UpdateProgress(w http.ResponseWriter, r *http.Request) {
	submissionID := chi.URLParam(r, "submission_id")
	var req struct {
		Stage           string `json:"stage"`
		ProgressPercent int    `json:"progress_percent"`
		CompanyID       string `json:"company_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", 400)
		return
	}

	if req.Stage != "" {
		if err := h.progressService.UpdateProgress(r.Context(), submissionID, req.Stage, req.ProgressPercent); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	}

	if err := h.progressService.LinkCompany(r.Context(), submissionID, req.CompanyID); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *FacultyIncubationHandler) GetCompanies(w http.ResponseWriter, r *http.Request) {
	companies, err := h.companyRepo.GetAll(r.Context())
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(companies)
}
