package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/rudraa2005/mic-website-main/backend/internal/middleware"
	"github.com/rudraa2005/mic-website-main/backend/internal/model"
	"github.com/rudraa2005/mic-website-main/backend/internal/service"
)

type SettingsHandler struct {
	settingService *service.SettingService
	profileService *service.ProfileService
	authService    *service.AuthService
}

type UpdateSettingsRequest struct {
	Profile struct {
		Name  string  `json:"name"`
		Phone *string `json:"phone,omitempty"`
		Bio   *string `json:"bio,omitempty"`
	} `json:"profile"`

	Preferences struct {
		Theme              string `json:"theme"`
		ApplicationUpdates bool   `json:"application_updates"`
		FeedbackAlerts     bool   `json:"feedback_alerts"`
		Newsletter         bool   `json:"newsletter"`
		EmailNotifications bool   `json:"email_notifications"`
	} `json:"preferences"`

	Security struct {
		Email *string `json:"email,omitempty"`
	} `json:"security"`
}

func NewSettingsHandler(s *service.SettingService, p *service.ProfileService, a *service.AuthService) *SettingsHandler {
	return &SettingsHandler{
		settingService: s,
		profileService: p,
		authService:    a,
	}
}

func (s *SettingsHandler) GetSettings(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, ok := middleware.GetUserFromContext(ctx)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	profile := &model.Profile{}
	p, err := s.profileService.GetProfile(ctx, user.UserID)
	if err == nil {
		profile = p
	}

	settings, err := s.settingService.GetSettings(ctx, user.UserID)
	if err != nil {
		http.Error(w, "settings not found", http.StatusNotFound)
		return
	}

	resp := map[string]interface{}{
		"profile": map[string]interface{}{
			"name":  profile.Name,
			"phone": profile.Phone,
			"bio":   profile.Bio,
			"email": profile.Email,
		},
		"preferences": map[string]interface{}{
			"theme":               settings.Theme,
			"application_updates": settings.ApplicationUpdates,
			"feedback_alerts":     settings.FeedbackAlerts,
			"newsletter":          settings.Newsletter,
			"email_notifications": settings.EmailNotifications,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (s *SettingsHandler) UpdateSettings(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, ok := middleware.GetUserFromContext(ctx)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req UpdateSettingsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if req.Profile.Name != "" || req.Profile.Phone != nil || req.Profile.Bio != nil {
		bio := ""
		if req.Profile.Bio != nil {
			bio = *req.Profile.Bio
		}

		if err := s.profileService.UpdateProfile(
			ctx,
			user.UserID,
			req.Profile.Name,
			req.Profile.Phone,
			bio,
		); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	settings := model.Settings{
		UserID:             user.UserID,
		Theme:              req.Preferences.Theme,
		ApplicationUpdates: req.Preferences.ApplicationUpdates,
		FeedbackAlerts:     req.Preferences.FeedbackAlerts,
		Newsletter:         req.Preferences.Newsletter,
		EmailNotifications: req.Preferences.EmailNotifications,
	}

	if err := s.settingService.UpdateSettings(ctx, settings); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("SETTINGS UPDATE: %+v\n", settings)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "settings saved",
	})
}
