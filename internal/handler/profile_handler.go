package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rudraa2005/mic-website-main/backend/internal/middleware"
	"github.com/rudraa2005/mic-website-main/backend/internal/service"
)

type ProfileHandler struct {
	profileService *service.ProfileService
}

func NewProfileHandler(s *service.ProfileService) *ProfileHandler {
	return &ProfileHandler{
		profileService: s,
	}
}

func (h *ProfileHandler) Me(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	profile, err := h.profileService.GetProfile(r.Context(), user.UserID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(profile)
}

func (h *ProfileHandler) UploadPhoto(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse multipart form with 10MB max file size
	err := r.ParseMultipartForm(10 << 20) // 10MB
	if err != nil {
		http.Error(w, "failed to parse form: "+err.Error(), http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("photo")
	if err != nil {
		http.Error(w, "failed to get file: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Validate file type
	ext := strings.ToLower(filepath.Ext(handler.Filename))
	allowedExts := []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}
	validExt := false
	for _, allowedExt := range allowedExts {
		if ext == allowedExt {
			validExt = true
			break
		}
	}
	if !validExt {
		http.Error(w, "invalid file type. Allowed: jpg, jpeg, png, gif, webp", http.StatusBadRequest)
		return
	}

	// Create uploads directory if it doesn't exist
	workDir, _ := os.Getwd()
	uploadsDir := filepath.Join(workDir, "frontend", "static", "uploads")
	if err := os.MkdirAll(uploadsDir, 0755); err != nil {
		http.Error(w, "failed to create uploads directory", http.StatusInternalServerError)
		return
	}

	// Generate unique filename
	filename := fmt.Sprintf("%s_%d%s", user.UserID, time.Now().Unix(), ext)
	filePath := filepath.Join(uploadsDir, filename)

	// Create file on disk
	dst, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "failed to save file: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Copy file content
	_, err = io.Copy(dst, file)
	if err != nil {
		os.Remove(filePath) // Clean up on error
		http.Error(w, "failed to save file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Generate URL path
	photoURL := fmt.Sprintf("/static/uploads/%s", filename)

	// Update profile with photo URL
	err = h.profileService.UpdatePhotoURL(r.Context(), user.UserID, photoURL)
	if err != nil {
		os.Remove(filePath) // Clean up on error
		http.Error(w, "failed to update profile: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message":   "Photo uploaded successfully",
		"photo_url": photoURL,
	})
}
