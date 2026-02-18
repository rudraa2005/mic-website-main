package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rudraa2005/mic-website-main/backend/internal/middleware"
	"github.com/rudraa2005/mic-website-main/backend/internal/repository"
)

type AdminSubmissionHandler struct {
	repo *repository.AdminSubmissionRepo
}

func NewAdminSubmissionHandler(repo *repository.AdminSubmissionRepo) *AdminSubmissionHandler {
	return &AdminSubmissionHandler{repo: repo}
}

// GetPendingSubmissions returns all submissions awaiting admin review
func (h *AdminSubmissionHandler) GetPendingSubmissions(w http.ResponseWriter, r *http.Request) {
	submissions, err := h.repo.GetPendingSubmissions(r.Context())
	if err != nil {
		log.Println("[ADMIN] GetPendingSubmissions failed:", err)
		http.Error(w, "failed to fetch submissions", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(submissions)
}

// DecideSubmission handles admin approve/reject decision
func (h *AdminSubmissionHandler) DecideSubmission(w http.ResponseWriter, r *http.Request) {
	submissionID := chi.URLParam(r, "id")

	var req struct {
		Decision string `json:"decision"` // "approved" or "rejected"
		Reason   string `json:"reason"`   // optional rejection reason
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("[ADMIN] DecideSubmission decode error:", err)
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

	var err error
	switch req.Decision {
	case "approved":
		err = h.repo.ApproveForFaculty(r.Context(), submissionID)
	case "rejected":
		err = h.repo.RejectSubmission(r.Context(), submissionID, req.Reason)
	default:
		http.Error(w, "invalid decision: must be 'approved' or 'rejected'", http.StatusBadRequest)
		return
	}

	if err != nil {
		log.Println("[ADMIN] DecideSubmission failed:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"success": true}`))
}

// GetAllSubmissions returns ALL submissions (not just pending) for admin view
func (h *AdminSubmissionHandler) GetAllSubmissions(w http.ResponseWriter, r *http.Request) {
	submissions, err := h.repo.GetAllSubmissions(r.Context())
	if err != nil {
		log.Println("[ADMIN] GetAllSubmissions failed:", err)
		http.Error(w, "failed to fetch submissions", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(submissions)
}

// AssignFaculty assigns a faculty member to a submission
func (h *AdminSubmissionHandler) AssignFaculty(w http.ResponseWriter, r *http.Request) {
	submissionID := chi.URLParam(r, "id")

	var req struct {
		FacultyID string `json:"faculty_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

	// Get admin user ID from context using the middleware helper
	claims, err := middleware.GetUser(r)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	adminID := claims.UserID

	err = h.repo.AssignFacultyToSubmission(r.Context(), submissionID, req.FacultyID, adminID)
	if err != nil {
		log.Println("[ADMIN] AssignFaculty failed:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"success": true}`))
}

// RemoveFaculty removes a faculty assignment from a submission
func (h *AdminSubmissionHandler) RemoveFaculty(w http.ResponseWriter, r *http.Request) {
	submissionID := chi.URLParam(r, "id")
	facultyID := chi.URLParam(r, "faculty_id")

	err := h.repo.RemoveFacultyFromSubmission(r.Context(), submissionID, facultyID)
	if err != nil {
		log.Println("[ADMIN] RemoveFaculty failed:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"success": true}`))
}

// GetAssignedFaculty returns faculty assigned to a submission
func (h *AdminSubmissionHandler) GetAssignedFaculty(w http.ResponseWriter, r *http.Request) {
	submissionID := chi.URLParam(r, "id")

	faculty, err := h.repo.GetAssignedFaculty(r.Context(), submissionID)
	if err != nil {
		log.Println("[ADMIN] GetAssignedFaculty failed:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(faculty)
}

// UpdateTags updates tags and domain for a submission
func (h *AdminSubmissionHandler) UpdateTags(w http.ResponseWriter, r *http.Request) {
	submissionID := chi.URLParam(r, "id")

	var req struct {
		Tags   []string `json:"tags"`
		Domain string   `json:"domain"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

	err := h.repo.UpdateSubmissionTags(r.Context(), submissionID, req.Tags, req.Domain)
	if err != nil {
		log.Println("[ADMIN] UpdateTags failed:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"success": true}`))
}
