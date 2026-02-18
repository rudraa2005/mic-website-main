package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rudraa2005/mic-website-main/backend/internal/middleware"
	"github.com/rudraa2005/mic-website-main/backend/internal/service"
)

type AIHandler struct {
	ai *service.AIService
}

func NewAIHandler(ai *service.AIService) *AIHandler {
	return &AIHandler{ai: ai}
}

func (h *AIHandler) AnalyzeDraft(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		SubmissionID string `json:"submission_id"`
		FilePath     string `json:"file_path"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	absPath := fmt.Sprintf("./uploads/%s_%s", req.SubmissionID, req.FilePath)
	log.Printf(
		"[AI] analyze request received | submission=%s user=%s file=%s\n",
		req.SubmissionID,
		user.UserID,
		absPath,
	)

	err := h.ai.CreateDraft(
		r.Context(),
		req.SubmissionID,
		user.UserID,
		absPath,
	)
	if err != nil {
		log.Println("[AI] CreateDraft failed:", err)
		http.Error(w, "failed to create ai draft", 500)
		return
	}
	log.Println("[AI] draft row created in DB")

	h.ai.AnalyzeDraft(r.Context(), req.SubmissionID, absPath)
	log.Println("[AI] analysis started (background)")

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "analysis_started",
	})
}

// frontend polls this
func (h *AIHandler) GetInsights(w http.ResponseWriter, r *http.Request) {
	_, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	submissionID := chi.URLParam(r, "submission_id")

	insights, err := h.ai.GetInsights(r.Context(), submissionID)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"status":   "completed",
		"insights": insights,
	})
}
