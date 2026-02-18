package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/rudraa2005/mic-website-main/backend/internal/middleware"
	"github.com/rudraa2005/mic-website-main/backend/internal/service"
)

type EventInvitationHandler struct {
	service *service.EventInvitationService
}

func NewEventInvitationHandler(
	service *service.EventInvitationService,
) *EventInvitationHandler {
	return &EventInvitationHandler{service: service}
}

/*
API RESPONSE DTO
Frontend depends on this shape
*/
type FacultyEventResponse struct {
	InvitationID string    `json:"invitation_id"`
	Title        string    `json:"title"`
	EventDate    string    `json:"date"`
	Venue        string    `json:"location"`
	Price        string    `json:"price"`
	Status       string    `json:"status"`
	InvitedAt    time.Time `json:"invited_at"`
}

/*
GET /faculty/events/invitations
*/
func (h *EventInvitationHandler) GetMyInvitations(
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

	events, err := h.service.GetFacultyEvents(
		r.Context(),
		claims.UserID,
	)
	if err != nil {
		http.Error(w, "failed to fetch invitations", http.StatusInternalServerError)
		return
	}

	// 🔑 MAP service models → API response
	resp := make([]FacultyEventResponse, 0, len(events))
	for _, e := range events {
		resp = append(resp, FacultyEventResponse{
			InvitationID: e.InvitationID,
			Title:        e.Title,
			EventDate:    e.EventDate,
			Venue:        e.Venue,
			Price:        e.Price,
			Status:       e.Status,
			InvitedAt:    e.InvitedAt,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

/*
POST /faculty/events/invitations/{invitation_id}/rsvp
*/
type RSVPRequest struct {
	Status string `json:"status"` // accepted | declined
}

func (h *EventInvitationHandler) UpdateRSVP(
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

	invitationID := chi.URLParam(r, "invitation_id")
	if invitationID == "" {
		http.Error(w, "missing invitation id", http.StatusBadRequest)
		return
	}

	var req RSVPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Status != "accepted" && req.Status != "declined" {
		http.Error(w, "invalid status", http.StatusBadRequest)
		return
	}

	err = h.service.UpdateRSVP(
		r.Context(),
		invitationID,
		claims.UserID,
		req.Status,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
