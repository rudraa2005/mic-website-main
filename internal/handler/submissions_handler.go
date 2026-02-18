package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rudraa2005/mic-website-main/backend/internal/middleware"
	"github.com/rudraa2005/mic-website-main/backend/internal/model"
	"github.com/rudraa2005/mic-website-main/backend/internal/service"
)

type SubmissionsHandler struct {
	submissionsService *service.SubmissionsService
}
type CreateSubmissionRequest struct {
	SubmissionID string  `json:"submission_id,omitempty"`
	Title        string  `json:"title"`
	Description  string  `json:"description"`
	FilePath     *string `json:"file_path"`
}

func NewSubmissionsHandler(ss *service.SubmissionsService) *SubmissionsHandler {
	return &SubmissionsHandler{
		submissionsService: ss,
	}
}

type SubmissionResponse struct {
	SubmissionID string    `json:"submission_id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	Status       string    `json:"status"`
	FilePath     *string   `json:"file_path"`
	CreatedAt    time.Time `json:"created_at"`
}
type SubmitSubmissionRequest struct {
	SubmissionID string `json:"submission_id"`
}

func (sh *SubmissionsHandler) CreateSubmission(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, ok := middleware.GetUserFromContext(ctx)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req CreateSubmissionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	submission := &model.Submission{
		SubmissionID: uuid.NewString(),
		UserID:       user.UserID,
		Title:        req.Title,
		Description:  req.Description,
		Status:       "draft",
	}

	if err := sh.submissionsService.Create(ctx, submission); err != nil {
		log.Println("CREATE SUBMISSION ERROR:", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"submission_id": submission.SubmissionID,
	})
}

func (sh *SubmissionsHandler) UpdateSubmission(w http.ResponseWriter, r *http.Request) {
	log.Println("UpdateSubmission HIT")
	ctx := r.Context()

	user, ok := middleware.GetUserFromContext(ctx)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	submissionID := chi.URLParam(r, "submission_id")
	if submissionID == "" {
		http.Error(w, "submission_id is required", http.StatusBadRequest)
		return
	}
	var req CreateSubmissionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	log.Println("UpdateSubmission: request decoded:", req)
	submission := &model.Submission{
		SubmissionID: submissionID,
		UserID:       user.UserID,
		Title:        req.Title,
		Description:  req.Description,
	}
	log.Println("UpdateSubmission: updating submission for userID:", user.UserID, submission)
	if err := sh.submissionsService.UpdateDraft(ctx, submission); err != nil {
		http.Error(w, "failed to update submission", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "file uploaded successfully",
	})
}

func (sh *SubmissionsHandler) GetByUserID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, ok := middleware.GetUserFromContext(ctx)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	submission, err := sh.submissionsService.GetByUserID(ctx, user.UserID)
	if err != nil {
		log.Println("Get submissions error:", err)
		http.Error(w, "failed to get submissions", http.StatusInternalServerError)
		return
	}
	var resp []SubmissionResponse
	for _, s := range submission {
		var filePath *string
		if s.FilePath != nil {
			filePath = s.FilePath
		}
		resp = append(resp, SubmissionResponse{
			SubmissionID: s.SubmissionID,
			Title:        s.Title,
			Description:  s.Description,
			Status:       s.Status,
			FilePath:     filePath,
			CreatedAt:    s.CreatedAt,
		})
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (sh *SubmissionsHandler) GetBySubmissionID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	_, ok := middleware.GetUserFromContext(ctx)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	submissions := chi.URLParam(r, "submission_id")

	submission, err := sh.submissionsService.GetBySubmissionID(ctx, submissions)
	if err != nil {
		http.Error(w, "failed to get submission", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(submission)
}

func (sh *SubmissionsHandler) DeleteSubmission(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, ok := middleware.GetUserFromContext(ctx)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	submissionID := r.URL.Query().Get("submission_id")

	err := sh.submissionsService.Delete(ctx, submissionID, user.UserID)
	if err != nil {
		http.Error(w, "Failed to delete submission", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (sh *SubmissionsHandler) SubmitSubmission(w http.ResponseWriter, r *http.Request) {

	log.Println("SubmitSubmission HIT")
	log.Println("URL:", r.URL.Path)
	log.Println("Param submission_id:", chi.URLParam(r, "submission_id"))

	ctx := r.Context()

	user, ok := middleware.GetUserFromContext(ctx)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	submissionID := chi.URLParam(r, "submission_id")

	userID := user.UserID
	log.Println("SubmitSubmission: userID =", userID, "submissionID =", submissionID)
	if submissionID == "" {
		http.Error(w, "submission_id is required", http.StatusBadRequest)
		return
	}

	err := sh.submissionsService.Submit(ctx, submissionID, userID, user.Email)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "submitted",
	})
}

func (sh *SubmissionsHandler) UploadSubmissionFile(w http.ResponseWriter, r *http.Request) {
	log.Println("UPLOAD route hit")
	ctx := r.Context()

	user, ok := middleware.GetUserFromContext(ctx)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	if user.UserID == "" {
		log.Println("UPLOAD: user ID missing in context")
		http.Error(w, "unauthorized: user ID missing", http.StatusUnauthorized)
		return
	}
	log.Println("UPLOAD: userID =", user.UserID)

	submissionID := chi.URLParam(r, "submission_id")
	log.Println("UPLOAD: submissionID =", submissionID)

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "invalid multipart form", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		log.Println("UPLOAD: form file missing:", err)
		http.Error(w, "file missing", http.StatusBadRequest)
		return
	}
	defer file.Close()
	log.Println("UPLOAD: received file:", header.Filename)

	os.MkdirAll("./uploads", os.ModePerm)
	path := fmt.Sprintf("./uploads/%s_%s", submissionID, header.Filename)

	dst, err := os.Create(path)
	if err != nil {
		log.Println("UPLOAD: failed to create file:", err)
		http.Error(w, "failed to save file", http.StatusInternalServerError)
		return
	}

	defer dst.Close()
	_, err = io.Copy(dst, file)
	if err != nil {
		log.Println("UPLOAD: failed to save file:", err)
		http.Error(w, "failed to save file", http.StatusInternalServerError)
		return
	}
	log.Println("UPLOAD: file saved to:", path)

	err = sh.submissionsService.AttachFile(ctx, submissionID, user.UserID, path)
	if err != nil {
		log.Println("UPLOAD: failed to attach file to submission:", err)
		http.Error(w, "failed to attach file", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"file_path": path,
	})
}
func (sh *SubmissionsHandler) DownloadSubmissionFile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	log.Println("DOWNLOAD route hit")
	user, ok := middleware.GetUserFromContext(ctx)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	submissionID := chi.URLParam(r, "submission_id")
	if submissionID == "" {
		http.Error(w, "submission_id missing", http.StatusBadRequest)
		return
	}

	submission, err := sh.submissionsService.GetBySubmissionID(ctx, submissionID)
	if err != nil {
		http.Error(w, "submission not found", http.StatusNotFound)
		return
	}

	// ownership check
	if submission.UserID != user.UserID && user.Role != "FACULTY" && user.Role != "ADMIN" {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	if submission.FilePath == nil {
		http.Error(w, "no file attached", http.StatusNotFound)
		return
	}

	http.ServeFile(w, r, *submission.FilePath)
}

func (sh *SubmissionsHandler) GetInsights(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	_, ok := middleware.GetUserFromContext(ctx)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	submissionID := chi.URLParam(r, "submission_id")

	insights, err := sh.submissionsService.GetAIInsights(ctx, submissionID)
	if err != nil {
		http.Error(w, "insights not ready", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"insights": insights,
	})
}
