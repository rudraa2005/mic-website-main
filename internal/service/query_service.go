package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/rudraa2005/mic-website-main/backend/internal/model"
	"github.com/rudraa2005/mic-website-main/backend/internal/repository"
)

type QueryService struct {
	queryRepo *repository.QueryRepo
}

func NewQueryService(
	queryRepo *repository.QueryRepo,
) *QueryService {
	return &QueryService{
		queryRepo: queryRepo,
	}
}

func (s *QueryService) CreateQuery(
	ctx context.Context,
	userID string,
	facultyID string,
	feedbackID *string,
	queryText string,
	priority string,
) error {

	if userID == "" {
		return errors.New("unauthorized")
	}

	if facultyID == "" {
		return errors.New("faculty_id required")
	}

	if queryText == "" {
		return errors.New("query cannot be empty")
	}

	if priority == "" {
		priority = "normal"
	}

	q := &model.Query{
		QueryID:    uuid.NewString(),
		UserID:     userID,
		FacultyID:  facultyID,
		FeedbackID: feedbackID,
		QueryText:  queryText,
		Priority:   priority,
		Status:     "pending",
	}

	return s.queryRepo.Create(ctx, q)
}

func (s *QueryService) GetMyQueries(
	ctx context.Context,
	userID string,
) ([]model.Query, error) {

	if userID == "" {
		return nil, errors.New("unauthorized")
	}

	return s.queryRepo.GetByUserID(ctx, userID)
}
