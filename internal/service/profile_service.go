package service

import (
	"context"
	"errors"

	"github.com/rudraa2005/mic-website-main/backend/internal/model"
)

type ProfileRepo interface {
	GetUserByID(ctx context.Context, userID string) (*model.Profile, error)
	UpdatePhotoURL(ctx context.Context, userID string, photoURL string) error
	UpdateProfile(ctx context.Context, userID string, name string, phone *string, bio string) error
}

type ProfileService struct {
	profileRepo ProfileRepo
}

func NewProfileService(profileRepo ProfileRepo) *ProfileService {
	return &ProfileService{
		profileRepo: profileRepo,
	}
}

func (ps *ProfileService) GetProfile(ctx context.Context, userID string) (*model.Profile, error) {
	return ps.profileRepo.GetUserByID(ctx, userID)
}

func (ps *ProfileService) UpdateProfile(ctx context.Context, userID string, name string, phone *string, bio string) error {

	profile, err := ps.profileRepo.GetUserByID(ctx, userID)
	if err != nil {
		return errors.New("Did not get user")
	}

	profile.Bio = bio
	profile.Phone = phone
	profile.Name = name

	err = ps.profileRepo.UpdateProfile(ctx, userID, name, phone, bio)
	return err
}

func (ps *ProfileService) UpdatePhotoURL(ctx context.Context, userID string, photoURL string) error {
	return ps.profileRepo.UpdatePhotoURL(ctx, userID, photoURL)
}
