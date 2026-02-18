package service

import (
	"context"

	"github.com/rudraa2005/mic-website-main/backend/internal/model"
)

type StartupRepository interface {
	Create(ctx context.Context, s *model.Startup) error
	GetByID(ctx context.Context, id string) (*model.Startup, error)
	ListByOwner(ctx context.Context, ownerID string) ([]model.Startup, error)
}

type StartupService struct {
	startup StartupRepository
}

func NewStartupService(startup StartupRepository) *StartupService {
	return &StartupService{
		startup: startup,
	}
}

func (s *StartupService) CreateStartup(ctx context.Context, startup *model.Startup) error {
	return s.startup.Create(ctx, startup)
}

func (s *StartupService) GetStartup(ctx context.Context, userID string) (*model.Startup, error) {
	return s.startup.GetByID(ctx, userID)
}

func (s *StartupService) ListMine(ctx context.Context, ownerID string) ([]model.Startup, error) {
	return s.startup.ListByOwner(ctx, ownerID)
}
