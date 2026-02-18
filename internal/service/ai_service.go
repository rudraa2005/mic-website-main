package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"path/filepath"

	"github.com/rudraa2005/mic-website-main/backend/internal/repository"
)

type AIService struct {
	baseURL string
	repo    *repository.AIRepo
	client  *http.Client
}

func NewAIService(
	baseURL string,
	aiRepo *repository.AIRepo,
	client *http.Client,
) *AIService {
	return &AIService{
		baseURL: baseURL,
		repo:    aiRepo,
		client:  client,
	}
}

func (s *AIService) AnalyzeDraft(
	ctx context.Context,
	submissionID string,
	absPath string,
) {
	// Python expects ONLY filename, not full path
	abs, err := filepath.Abs(absPath)
	if err != nil {
		log.Println("AI abs path erro:", err)
		return
	}

	payload := map[string]string{
		"submission_id": submissionID,
		"file_path":     abs,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		log.Println("AI marshal payload failed:", err)
		return
	}

	req, err := http.NewRequest(
		"POST",
		s.baseURL+"/analyze",
		bytes.NewBuffer(body),
	)
	if err != nil {
		log.Println("AI request creation failed:", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	log.Println("[AI -> PYTHON]POST", s.baseURL+"/analyze")
	log.Printf("[AI → PYTHON] Payload: submission_id=%s file_path=%s\n",
		submissionID,
		filepath.Base(absPath),
	)
	log.Println("[AI → PYTHON] Content-Type:", req.Header.Get("Content-Type"))
	resp, err := s.client.Do(req)
	if err != nil {
		log.Println("AI service call failed:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		log.Println("AI returned unexpected status:", resp.Status)
		return
	}

	var result map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Println("AI decode failed:", err)
		return
	}

	raw, _ := json.Marshal(result)

	if err := s.repo.SaveInsights(ctx, submissionID, string(raw)); err != nil {
		log.Println("AI save insights failed:", err)
		return
	}

	log.Println("AI analysis completed for submission:", submissionID)
}

func (s *AIService) CreateDraft(
	ctx context.Context,
	submissionID string,
	userID string,
	filePath string,
) error {
	return s.repo.CreateDraft(ctx, submissionID, userID, filePath)
}

func (s *AIService) GetDraftInsights(
	ctx context.Context,
	submissionID string,
) (string, string, error) {
	return s.repo.GetInsights(ctx, submissionID)
}

func (s *AIService) GetInsights(
	ctx context.Context,
	submissionID string,
) (string, error) {
	insights, status, err := s.GetDraftInsights(ctx, submissionID)
	if err != nil {
		return "", err
	}

	if status != "completed" {
		return "", errors.New("analysis not ready")
	}

	return insights, nil
}
