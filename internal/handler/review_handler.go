package handler

import (
	"encoding/json"
	"net/http"
	"github.com/go-chi/chi/v5"
	"github.com/rudraa2005/mic-website-main/backend/internal/service"
	"github.com/rudraa2005/mic-website-main/backend/internal/repository"
)

type ReviewHandler struct {
    Svc *service.ReviewService
}

func NewReviewHandler(s *service.ReviewService) *ReviewHandler {
    return &ReviewHandler{Svc: s}
}

func (h *ReviewHandler) SubmitReview(w http.ResponseWriter, r *http.Request) {
    var rev repository.Review
    json.NewDecoder(r.Body).Decode(&rev)

    rev.ReviewerID = r.Context().Value("userID").(string)
    rev.ReviewerName = r.Context().Value("name").(string)
    rev.ReviewerDesignation = r.Context().Value("designation").(string)

    err := h.Svc.SubmitReview(r.Context(), rev)

    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
}

func (h *ReviewHandler) ListForStartup(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    reviews, _ := h.Svc.GetStartupReviews(r.Context(), id)
    json.NewEncoder(w).Encode(reviews)
}

func (h *ReviewHandler) ListMyReviews(w http.ResponseWriter, r *http.Request) {
    reviewerID := r.Context().Value("userID").(string)
    list, _ := h.Svc.GetMyReviews(r.Context(), reviewerID)
    json.NewEncoder(w).Encode(list)
}
