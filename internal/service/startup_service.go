package service

import(
	"context"
	"github.com/rudraa2005/mic-website-main/backend/internal/model"
    "github.com/rudraa2005/mic-website-main/backend/internal/repository"
)

type StartupService struct{
	repo *repository.StartupRepository
}

func NewStartupService( repo *repository.StartupRepository) *StartupService{
	return &StartupService{repo}
}

func (s *StartupService) CreateStartup(ctx context.Context, startup *model.Startup) error {
	return s.repo.Create(ctx, startup)
}

func (s *StartupService) GetStartup(ctx context.Context, id string) (*model.Startup, error) {
    return s.repo.GetByID(ctx, id)
}

func (s *StartupService) ListMine(ctx context.Context, ownerID string) ([]model.Startup, error) {
    return s.repo.ListByOwner(ctx, ownerID)
}