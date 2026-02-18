package service

import (
	"context"
	"log"

	"github.com/rudraa2005/mic-website-main/backend/internal/model"
)

type SubmissionsRepo interface {
	Create(ctx context.Context, s *model.Submission) error
	UpdateDraft(ctx context.Context, s *model.Submission) error
	GetByUserID(ctx context.Context, userID string) ([]model.Submission, error)
	GetBySubmissionID(ctx context.Context, submissionID string) (*model.Submission, error)
	Delete(ctx context.Context, submissionID string, userID string) error
	MarkSubmitted(ctx context.Context, submissionID string, userID string) error
	AttachFile(ctx context.Context, submissionID string, userID string, filePath string) error

	UpdateStatus(
		ctx context.Context,
		submissionID string,
		newStatus string,
		newStage string,
	) error
}

type SubmissionsService struct {
	NotificationService *NotificationService
	submissionsRepo     SubmissionsRepo
	UserRepo            ProfileRepo
	AIService           *AIService
}

func NewSubmissionsService(
	notificationService *NotificationService,
	submissionsRepo SubmissionsRepo,
	UserRepo ProfileRepo,
	AIService *AIService,
) *SubmissionsService {
	return &SubmissionsService{
		NotificationService: notificationService,
		submissionsRepo:     submissionsRepo,
		UserRepo:            UserRepo,
		AIService:           AIService,
	}
}

func (s *SubmissionsService) Create(ctx context.Context, submission *model.Submission) error {
	return s.submissionsRepo.Create(ctx, submission)
}

func (s *SubmissionsService) UpdateDraft(ctx context.Context, submission *model.Submission) error {
	return s.submissionsRepo.UpdateDraft(ctx, submission)
}

func (s *SubmissionsService) GetByUserID(
	ctx context.Context,
	userID string,
) ([]model.Submission, error) {
	submission, _ := s.submissionsRepo.GetByUserID(ctx, userID)
	log.Println("SERVICE: fetched", len(submission), "submissions for userID:", userID, submission)
	return submission, nil
}

func (s *SubmissionsService) GetBySubmissionID(ctx context.Context, submissionID string) (*model.Submission, error) {
	return s.submissionsRepo.GetBySubmissionID(ctx, submissionID)
}

func (s *SubmissionsService) Delete(ctx context.Context, submissionID string, userID string) error {
	return s.submissionsRepo.Delete(ctx, submissionID, userID)
}

func (s *SubmissionsService) Submit(
	ctx context.Context,
	submissionID string,
	userID string,
	email string,
) error {

	if err := s.submissionsRepo.MarkSubmitted(ctx, submissionID, userID); err != nil {
		return err
	}

	submission, err := s.submissionsRepo.GetBySubmissionID(ctx, submissionID)
	if err != nil {
		return err
	}
	user, err := s.UserRepo.GetUserByID(ctx, userID)
	if err != nil || len(user.UserID) == 0 {
		return err
	}

	log.Println("[SUBMIT SERVICE] status updated, calling notification")
	err = s.NotificationService.NotifyStatusChange(
		ctx,
		submission.UserID,
		email,
		submissionID,
		"draft",
		"submitted",
	)
	if err != nil {
		log.Println("notification failed:", err)
	}
	return nil
}

func (s *SubmissionsService) AttachFile(
	ctx context.Context,
	submissionID string,
	userID string,
	filePath string,
) error {
	return s.submissionsRepo.AttachFile(ctx, submissionID, userID, filePath)
}

func (s *SubmissionsService) UpdateStatus(
	ctx context.Context,
	submissionID string,
	newStatus string,
	user model.User,
) error {

	submission, err := s.submissionsRepo.GetBySubmissionID(ctx, submissionID)
	if err != nil {
		return err
	}

	oldStatus := submission.Status
	if oldStatus == newStatus {
		return nil
	}

	err = s.submissionsRepo.UpdateStatus(
		ctx,
		submissionID,
		newStatus,
		submission.Stage,
	)
	if err != nil {
		log.Println("notification failed:", err)
	}

	go s.NotificationService.NotifyStatusChange(
		ctx,
		user.UserID,
		user.Email,
		submissionID,
		oldStatus,
		newStatus,
	)

	return nil
}

func (s *SubmissionsService) GetAIInsights(
	ctx context.Context,
	submissionID string,
) (string, error) {
	return s.AIService.GetInsights(ctx, submissionID)
}
