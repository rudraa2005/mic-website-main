package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rudraa2005/mic-website-main/backend/internal/repository"
)

type AdminWorkHandler struct {
	repo *repository.AdminWorkRepo
}

func NewAdminWorkHandler(repo *repository.AdminWorkRepo) *AdminWorkHandler {
	return &AdminWorkHandler{repo: repo}
}

// GetAllWork returns all work items for admin management
func (h *AdminWorkHandler) GetAllWork(w http.ResponseWriter, r *http.Request) {
	work, err := h.repo.GetAllWork(r.Context())
	if err != nil {
		log.Println("[ADMIN] GetAllWork failed:", err)
		http.Error(w, "failed to fetch work items", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(work)
}

// UpdateWork updates a work item
func (h *AdminWorkHandler) UpdateWork(w http.ResponseWriter, r *http.Request) {
	workID := chi.URLParam(r, "id")

	var req struct {
		Title           string  `json:"title"`
		Description     string  `json:"description"`
		Stage           string  `json:"stage"`
		ProgressPercent int     `json:"progress_percent"`
		CompanyID       *string `json:"company_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

	err := h.repo.UpdateWork(r.Context(), workID, req.Title, req.Description, req.Stage, req.ProgressPercent, req.CompanyID)
	if err != nil {
		log.Println("[ADMIN] UpdateWork failed:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"success": true}`))
}

// DeleteWork removes a work item
func (h *AdminWorkHandler) DeleteWork(w http.ResponseWriter, r *http.Request) {
	workID := chi.URLParam(r, "id")

	err := h.repo.DeleteWork(r.Context(), workID)
	if err != nil {
		log.Println("[ADMIN] DeleteWork failed:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"success": true}`))
}

// GetCompanies returns all companies
func (h *AdminWorkHandler) GetCompanies(w http.ResponseWriter, r *http.Request) {
	companies, err := h.repo.GetCompanies(r.Context())
	if err != nil {
		log.Println("[ADMIN] GetCompanies failed:", err)
		http.Error(w, "failed to fetch companies", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(companies)
}

// AddCompany creates a new company
func (h *AdminWorkHandler) AddCompany(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name    string  `json:"name"`
		LogoURL *string `json:"logo_url"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

	id, err := h.repo.AddCompany(r.Context(), req.Name, req.LogoURL)
	if err != nil {
		log.Println("[ADMIN] AddCompany failed:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"id": id})
}

// DeleteCompany removes a company
func (h *AdminWorkHandler) DeleteCompany(w http.ResponseWriter, r *http.Request) {
	companyID := chi.URLParam(r, "id")

	err := h.repo.DeleteCompany(r.Context(), companyID)
	if err != nil {
		log.Println("[ADMIN] DeleteCompany failed:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"success": true}`))
}
