package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rudraa2005/mic-website-main/backend/internal/repository"
	"github.com/rudraa2005/mic-website-main/backend/internal/service"
)

type ContentHandler struct {
	service *service.ContentService
}

func NewContentHandler(service *service.ContentService) *ContentHandler {
	return &ContentHandler{service: service}
}

/* -------------------- ABOUT -------------------- */

func (h *ContentHandler) GetAboutCards(w http.ResponseWriter, r *http.Request) {
	items, err := h.service.GetAboutCards(r.Context())
	if err != nil {
		log.Println("[CONTENT] GetAboutCards failed:", err)
		http.Error(w, "failed to fetch about cards", http.StatusInternalServerError)
		return
	}

	writeJSON(w, items)
}

func (h *ContentHandler) GetAboutFeatures(w http.ResponseWriter, r *http.Request) {
	items, err := h.service.GetAboutFeatures(r.Context())
	if err != nil {
		log.Println("[CONTENT] GetAboutFeatures failed:", err)
		http.Error(w, "failed to fetch about features", http.StatusInternalServerError)
		return
	}

	writeJSON(w, items)
}

func (h *ContentHandler) GetTeamMembers(w http.ResponseWriter, r *http.Request) {
	items, err := h.service.GetTeamMembers(r.Context())
	if err != nil {
		log.Println("[CONTENT] GetTeamMembers failed:", err)
		http.Error(w, "failed to fetch team members", http.StatusInternalServerError)
		return
	}
	writeJSON(w, items)
}

func (h *ContentHandler) GetTestimonials(w http.ResponseWriter, r *http.Request) {
	items, err := h.service.GetTestimonials(r.Context())
	if err != nil {
		log.Println("[CONTENT] GetTestimonials failed:", err)
		http.Error(w, "failed to fetch testimonials", http.StatusInternalServerError)
		return
	}
	writeJSON(w, items)
}

func (h *ContentHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	items, err := h.service.GetStats(r.Context())
	if err != nil {
		log.Println("[CONTENT] GetStats failed:", err)
		http.Error(w, "failed to fetch stats", http.StatusInternalServerError)
		return
	}
	writeJSON(w, items)
}

/* -------------------- RESOURCES -------------------- */

func (h *ContentHandler) GetResources(w http.ResponseWriter, r *http.Request) {
	items, err := h.service.GetResources(r.Context())
	if err != nil {
		log.Println("[CONTENT] GetResources failed:", err)
		http.Error(w, "failed to fetch resources", http.StatusInternalServerError)
		return
	}

	writeJSON(w, items)
}

func (h *ContentHandler) GetTopResources(w http.ResponseWriter, r *http.Request) {
	items, err := h.service.GetTopResources(r.Context())
	if err != nil {
		log.Println("[CONTENT] GetTopResources failed:", err)
		http.Error(w, "failed to load top resources", http.StatusInternalServerError)
		return
	}

	writeJSON(w, items)
}

/* -------------------- EVENTS -------------------- */

func (h *ContentHandler) GetUpcomingEvents(w http.ResponseWriter, r *http.Request) {
	events, err := h.service.GetUpcomingEvents(r.Context())
	if err != nil {
		log.Println("[EVENTS] GetUpcomingEvents failed:", err)
		http.Error(w, "failed to fetch events", http.StatusInternalServerError)
		return
	}

	writeJSON(w, events)
}

func (h *ContentHandler) GetAllEvents(w http.ResponseWriter, r *http.Request) {
	events, err := h.service.GetAllEvents(r.Context())
	if err != nil {
		log.Println("[EVENTS] GetAllEvents failed:", err)
		http.Error(w, "failed to fetch events", http.StatusInternalServerError)
		return
	}

	writeJSON(w, events)
}

/* -------------------- helpers -------------------- */

func writeJSON(w http.ResponseWriter, payload any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(payload)
}

func (h *ContentHandler) GetAllContent(w http.ResponseWriter, r *http.Request) {
	items, err := h.service.GetAllContent(r.Context())
	if err != nil {
		http.Error(w, "failed to fetch content", 500)
		return
	}
	writeJSON(w, items)
}

type CreateContentRequest struct {
	ContentType string                 `json:"content_type"`
	Title       string                 `json:"title"`
	Description *string                `json:"description"`
	ImageURL    *string                `json:"image_url"`
	OrderIndex  int                    `json:"order_index"`
	IsActive    bool                   `json:"is_active"`
	ContentData map[string]interface{} `json:"content_data"`
}

func (h *ContentHandler) CreateContent(w http.ResponseWriter, r *http.Request) {
	var req CreateContentRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println(err)
		http.Error(w, "invalid payload", 400)
		return
	}

	// Serialize content_data to JSON bytes
	contentDataBytes := []byte(`{}`)
	if req.ContentData != nil {
		if data, err := json.Marshal(req.ContentData); err == nil {
			contentDataBytes = data
		}
	}

	c := &repository.Content{
		ContentType: req.ContentType,
		Title:       req.Title,
		Description: req.Description,
		ImageURL:    req.ImageURL,
		OrderIndex:  req.OrderIndex,
		IsActive:    req.IsActive,
		ContentData: contentDataBytes,
	}

	if err := h.service.CreateContent(r.Context(), c); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 500)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

type UpdateContentRequest struct {
	Title       string                 `json:"title"`
	Description *string                `json:"description"`
	ImageURL    *string                `json:"image_url"`
	OrderIndex  int                    `json:"order_index"`
	IsActive    bool                   `json:"is_active"`
	ContentData map[string]interface{} `json:"content_data"`
}

func (h *ContentHandler) UpdateContent(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	uid, _ := uuid.Parse(id)

	var req UpdateContentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid payload", 400)
		return
	}

	// Serialize content_data to JSON bytes
	contentDataBytes := []byte(`{}`)
	if req.ContentData != nil {
		if data, err := json.Marshal(req.ContentData); err == nil {
			contentDataBytes = data
		}
	}

	c := &repository.Content{
		Title:       req.Title,
		Description: req.Description,
		ImageURL:    req.ImageURL,
		OrderIndex:  req.OrderIndex,
		IsActive:    req.IsActive,
		ContentData: contentDataBytes,
	}

	if err := h.service.UpdateContent(r.Context(), uid, c); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ContentHandler) DeleteContent(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	uid, _ := uuid.Parse(id)

	if err := h.service.DeleteContent(r.Context(), uid); err != nil {
		http.Error(w, "failed to delete content", 500)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
