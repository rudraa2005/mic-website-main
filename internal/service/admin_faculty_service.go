package service

import (
	"context"
	"errors"

	"github.com/rudraa2005/mic-website-main/backend/internal/auth"
	"github.com/rudraa2005/mic-website-main/backend/internal/repository"
)

type AdminFacultyService struct {
	repo *repository.AdminFacultyRepository
}

func NewAdminFacultyService(repo *repository.AdminFacultyRepository) *AdminFacultyService {
	return &AdminFacultyService{repo: repo}
}

type FacultyResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (s *AdminFacultyService) GetAllFaculty(ctx context.Context) ([]FacultyResponse, error) {
	users, err := s.repo.GetAllByRole(ctx, "FACULTY")
	if err != nil {
		return nil, err
	}

	result := make([]FacultyResponse, len(users))
	for i, u := range users {
		result[i] = FacultyResponse{
			ID:    u.ID,
			Name:  u.Name,
			Email: u.Email,
		}
	}

	return result, nil
}

func (s *AdminFacultyService) CreateFaculty(ctx context.Context, name, email, password string) error {
	if name == "" || email == "" || password == "" {
		return errors.New("name, email and password are required")
	}

	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		return errors.New("failed to hash password")
	}

	return s.repo.Create(ctx, name, email, hashedPassword, "FACULTY")
}

func (s *AdminFacultyService) UpdateFaculty(ctx context.Context, id, name, email, password string) error {
	if id == "" {
		return errors.New("faculty id is required")
	}

	var hashedPassword string
	if password != "" {
		var err error
		hashedPassword, err = auth.HashPassword(password)
		if err != nil {
			return errors.New("failed to hash password")
		}
	}

	return s.repo.Update(ctx, id, name, email, hashedPassword)
}

func (s *AdminFacultyService) DeleteFaculty(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("faculty id is required")
	}
	return s.repo.Delete(ctx, id)
}
