package handler

import (
	"encoding/json"
	"net/http"

	"github.com/rudraa2005/mic-website-main/backend/internal/middleware"
	"github.com/rudraa2005/mic-website-main/backend/internal/model"
	"github.com/rudraa2005/mic-website-main/backend/internal/service"
)

type FeedbackHandler struct {
	feedbackService *service.FeedbackService
}

func NewFeedbackHandler(
	feedbackService *service.FeedbackService,
) *FeedbackHandler {
	return &FeedbackHandler{
		feedbackService: feedbackService,
	}
}

func (h *FeedbackHandler) GetMyFeedbacks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, ok := middleware.GetUserFromContext(ctx)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	feedbacks, err := h.feedbackService.GetMyFeedbacks(ctx, user.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(feedbacks)
}
func (h *FeedbackHandler) Create(w http.ResponseWriter, r *http.Request) {
	var f model.Feedback
	if err := json.NewDecoder(r.Body).Decode(&f); err != nil {
		http.Error(w, "invalid request", 400)
		return
	}

	claims, err := middleware.GetUser(r)
	if err != nil {
		http.Error(w, "unauthorized", 401)
		return
	}

	f.FacultyID = claims.UserID
	// We could fetch faculty details here, but for now we'll use what's sent or defaults
	if f.FacultyName == "" {
		f.FacultyName = claims.Email // Fallback to Email if Name is missing
	}
	f.Status = "active"

	if err := h.feedbackService.CreateFeedback(r.Context(), &f); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(f)
}
