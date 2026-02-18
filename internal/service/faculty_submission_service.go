package service

import (
	"context"
	"errors"

	"github.com/rudraa2005/mic-website-main/backend/internal/repository"
)

var (
	ErrInvalidDecision        = errors.New("invalid decision")
	ErrSubmissionNotProcessed = errors.New("submission could not be processed")
)

type FacultyReviewService struct {
	repo                *repository.FacultySubmissionRepo
	notificationService *NotificationService
}

func NewFacultyReviewService(
	repo *repository.FacultySubmissionRepo,
	notificationService *NotificationService,
) *FacultyReviewService {
	return &FacultyReviewService{
		repo:                repo,
		notificationService: notificationService,
	}
}

func (s *FacultyReviewService) GetSubmitted(
	ctx context.Context,
) ([]repository.FacultySubmission, error) {
	return s.repo.GetSubmitted(ctx)
}

func (s *FacultyReviewService) GetByID(
	ctx context.Context,
	id string,
) (*repository.FacultySubmission, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *FacultyReviewService) Decide(
	ctx context.Context,
	submissionID string,
	facultyID string,
	decision string,
) error {

	if submissionID == "" || facultyID == "" {
		return ErrSubmissionNotProcessed
	}

	switch decision {
	case "approved":
		email, title, err := s.repo.ApproveSubmission(ctx, submissionID, facultyID)
		if err == nil {
			s.notificationService.SendSubmissionStatusUpdate(ctx, email, title, "approved (Under Incubation)")
		}
		return err

	case "rejected":
		email, title, err := s.repo.RejectSubmission(ctx, submissionID, facultyID)
		if err == nil {
			s.notificationService.SendSubmissionStatusUpdate(ctx, email, title, "rejected")
		}
		return err

	case "needs_improvement":
		email, title, err := s.repo.NeedsImprovement(ctx, submissionID, facultyID)
		if err == nil {
			s.notificationService.SendSubmissionStatusUpdate(ctx, email, title, "needs improvement - please revise your submission")
		}
		return err

	default:
		return ErrInvalidDecision
	}
}
