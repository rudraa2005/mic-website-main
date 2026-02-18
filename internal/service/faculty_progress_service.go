package service

import (
	"context"

	"github.com/rudraa2005/mic-website-main/backend/internal/repository"
)

type FacultyProgressService struct {
	repo *repository.FacultyProgressRepository
}

func NewFacultyProgressService(
	repo *repository.FacultyProgressRepository,
) *FacultyProgressService {
	return &FacultyProgressService{repo: repo}
}

func (s *FacultyProgressService) GetProgressForFaculty(
	ctx context.Context,
	facultyID string,
) ([]repository.FacultyProgress, error) {
	return s.repo.GetByFaculty(ctx, facultyID)
}
func (s *FacultyProgressService) GetProgressBySubmission(
	ctx context.Context,
	facultyID string,
	submissionID string,
) (*repository.FacultyProgress, error) {
	return s.repo.GetBySubmission(ctx, facultyID, submissionID)
}
func (s *FacultyProgressService) UpdateProgress(ctx context.Context, submissionID string, stage string, progress int) error {
	return s.repo.UpdateProgress(ctx, submissionID, stage, progress)
}

func (s *FacultyProgressService) LinkCompany(ctx context.Context, submissionID string, companyID string) error {
	return s.repo.LinkCompany(ctx, submissionID, companyID)
}
