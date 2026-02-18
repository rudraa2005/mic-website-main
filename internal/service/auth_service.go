package service

import (
	"context"
	"errors"
	"strings"

	"github.com/rudraa2005/mic-website-main/backend/internal/auth"
	"github.com/rudraa2005/mic-website-main/backend/internal/model"
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrInvalidEmail      = errors.New("invalid email domain")
	ErrUserDoesNotExist  = errors.New("User does not exist, Signup first")
	ErrIncorrectPassword = errors.New("Incorrect Password")
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) (string, error)
	FindByEmail(email string) (*model.User, error)
	UpdatePassword(ctx context.Context, userID string, password string) error
	GetByID(ctx context.Context, userID string) (*model.User, error)
}

type FacultyRepository interface {
	FindByEmail(email string) (*model.Faculty, error)
}

type SettingsRepository interface {
	CreateDefaults(ctx context.Context, userID string) error
}
type ProfileRepository interface {
	StoreUser(ctx context.Context, p model.Profile) error
}

type AuthService struct {
	userRepo     UserRepository
	facultyRepo  FacultyRepository
	profileRepo  ProfileRepository
	settingsRepo SettingsRepository
}

func NewAuthService(
	userRepo UserRepository,
	facultyRepo FacultyRepository,
	profileRepo ProfileRepository,
	settingsRepo SettingsRepository,
) *AuthService {
	return &AuthService{
		userRepo:     userRepo,
		facultyRepo:  facultyRepo,
		profileRepo:  profileRepo,
		settingsRepo: settingsRepo,
	}
}

func (s *AuthService) Signup(ctx context.Context, email string, password string, name string) error {
	_, err := s.userRepo.FindByEmail(email)
	if err == nil {
		return ErrUserAlreadyExists
	}

	role, err := defineRole(email)
	if err != nil {
		return err
	}

	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		return err
	}

	user := &model.User{
		Email:          email,
		HashedPassword: hashedPassword,
		Role:           role,
		Name:           name,
	}

	if _, err := s.userRepo.Create(ctx, user); err != nil {
		return err
	}

	profile := model.Profile{
		UserID: user.UserID,
		Email:  email,
		Name:   name,
	}

	if err := s.profileRepo.StoreUser(ctx, profile); err != nil {
		return err
	}
	_ = s.settingsRepo.CreateDefaults(ctx, user.UserID)

	return nil
}

func defineRole(email string) (string, error) {
	if strings.HasSuffix(email, "@learner.manipal.edu") {
		return "STUDENT", nil
	}
	if strings.HasSuffix(email, "@manipal.edu") {
		return "FACULTY", nil
	}

	return "", ErrInvalidEmail
}

func (s *AuthService) Login(email string, password string) (*model.User, string, error) {

	user, err := s.userRepo.FindByEmail(email)
	if err == nil {
		if !auth.CompareHashedPassword(user.HashedPassword, password) {
			return nil, "", ErrIncorrectPassword
		}

		token, err := auth.CreateToken(user.UserID, user.Role, email)
		if err != nil {
			return nil, "", err
		}

		return user, token, nil
	}

	faculty, err := s.facultyRepo.FindByEmail(email)
	if err != nil {
		return nil, "", ErrUserDoesNotExist
	}

	if !auth.CompareHashedPassword(faculty.PasswordHash, password) {
		return nil, "", ErrIncorrectPassword
	}

	token, err := auth.CreateToken(faculty.ID, faculty.Role, faculty.Email)
	if err != nil {
		return nil, "", err
	}

	return &model.User{
		UserID: faculty.ID,
		Email:  faculty.Email,
		Role:   faculty.Role,
		Name:   faculty.Name,
	}, token, nil
}

func (s *AuthService) ChangePassword(ctx context.Context, userID string, currentPassword string, newPassword string) error {

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return ErrUserDoesNotExist
	}
	ok := auth.CompareHashedPassword(user.HashedPassword, currentPassword)
	newHashedPassword, err := auth.HashPassword(newPassword)
	if err != nil {
		return errors.New("Unable to Hash Password")
	}

	if !ok {
		return ErrIncorrectPassword
	}

	err = s.userRepo.UpdatePassword(ctx, userID, newHashedPassword)

	return err
}
