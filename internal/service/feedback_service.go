package service

import (
	"context"
	"errors"

	"github.com/rudraa2005/mic-website-main/backend/internal/model"
	"github.com/rudraa2005/mic-website-main/backend/internal/repository"
)

type FeedbackService struct {
	feedbackRepo *repository.FeedbackRepo
}

func NewFeedbackService(
	feedbackRepo *repository.FeedbackRepo,
) *FeedbackService {
	return &FeedbackService{
		feedbackRepo: feedbackRepo,
	}
}

func (s *FeedbackService) GetMyFeedbacks(
	ctx context.Context,
	userID string,
) ([]model.Feedback, error) {

	if userID == "" {
		return nil, errors.New("user id required")
	}

	return s.feedbackRepo.GetByUserID(ctx, userID)
}
func (s *FeedbackService) CreateFeedback(ctx context.Context, feedback *model.Feedback) error {
	return s.feedbackRepo.Create(ctx, feedback)
}
