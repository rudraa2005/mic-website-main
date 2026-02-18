package service

import (
	"context"

	"github.com/rudraa2005/mic-website-main/backend/internal/model"
)

type SettingsRepo interface {
	GetByUserID(ctx context.Context, userID string) (*model.Settings, error)
	Update(ctx context.Context, s model.Settings) error
	CreateDefaults(ctx context.Context, userID string) error
}

type SettingService struct {
	settingsRepo SettingsRepo
}

func NewSettingService(sr SettingsRepo) *SettingService {
	return &SettingService{
		settingsRepo: sr,
	}
}
func (s *SettingService) GetSettings(ctx context.Context, userID string) (*model.Settings, error) {
	settings, err := s.settingsRepo.GetByUserID(ctx, userID)
	if err != nil {
		if err := s.settingsRepo.CreateDefaults(ctx, userID); err != nil {
			return nil, err
		}
		return s.settingsRepo.GetByUserID(ctx, userID)
	}
	return settings, nil
}

func (s *SettingService) UpdateSettings(ctx context.Context, ms model.Settings) error {
	return s.settingsRepo.Update(ctx, ms)
}

func (s *SettingService) CreateDefaults(ctx context.Context, userID string) error {
	return s.settingsRepo.CreateDefaults(ctx, userID)
}
