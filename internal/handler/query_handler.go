package handler

import (
	"encoding/json"
	"net/http"

	"github.com/rudraa2005/mic-website-main/backend/internal/middleware"
	"github.com/rudraa2005/mic-website-main/backend/internal/service"
)

type QueryHandler struct {
	queryService *service.QueryService
}

func NewQueryHandler(
	queryService *service.QueryService,
) *QueryHandler {
	return &QueryHandler{
		queryService: queryService,
	}
}

type CreateQueryRequest struct {
	FacultyID  string  `json:"faculty_id"`
	FeedbackID *string `json:"feedback_id"`
	Query      string  `json:"query"`
	Priority   string  `json:"priority"`
}

func (h *QueryHandler) CreateQuery(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, ok := middleware.GetUserFromContext(ctx)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req CreateQueryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	err := h.queryService.CreateQuery(
		ctx,
		user.UserID,
		req.FacultyID,
		req.FeedbackID,
		req.Query,
		req.Priority,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *QueryHandler) GetMyQueries(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, ok := middleware.GetUserFromContext(ctx)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	queries, err := h.queryService.GetMyQueries(ctx, user.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(queries)
}
