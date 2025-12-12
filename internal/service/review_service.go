package service

import (
    "context"
    "github.com/rudraa2005/mic-website-main/backend/internal/repository"
)

type ReviewService struct {
    Repo *repository.ReviewRepository
}

func NewReviewService(repo *repository.ReviewRepository) *ReviewService {
    return &ReviewService{Repo: repo}
}

func (s *ReviewService) SubmitReview(ctx context.Context, r repository.Review) error {
    return s.Repo.CreateReview(ctx, r)
}

func (s *ReviewService) GetStartupReviews(ctx context.Context, startupID string) ([]repository.Review, error) {
    return s.Repo.ListForStartup(ctx, startupID)
}

func (s *ReviewService) GetMyReviews(ctx context.Context, reviewerID string) ([]repository.Review, error) {
    return s.Repo.ListForReviewer(ctx, reviewerID)
}
