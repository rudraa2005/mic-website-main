package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rudraa2005/mic-website-main/backend/internal/model"
	"github.com/rudraa2005/mic-website-main/backend/internal/repository"
)

type ContentService struct {
	repo *repository.ContentRepository
}

func NewContentService(repo *repository.ContentRepository) *ContentService {
	return &ContentService{repo: repo}
}

func derefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func mapBaseContent(c repository.Content) model.BaseContent {
	return model.BaseContent{
		ID:        c.ID.String(),
		Title:     c.Title,
		IsActive:  c.IsActive,
		Order:     c.OrderIndex,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}

func mapToAbout(items []repository.Content) []model.AboutContent {
	out := make([]model.AboutContent, 0, len(items))

	for _, c := range items {
		out = append(out, model.AboutContent{
			BaseContent: mapBaseContent(c),
			Description: derefString(c.Description),
		})
	}

	return out
}

func (s *ContentService) GetAboutCards(ctx context.Context) ([]model.AboutContent, error) {
	items, err := s.repo.GetActiveByType(ctx, "about_card")
	if err != nil {
		return nil, err
	}
	return mapToAbout(items), nil
}

func (s *ContentService) GetAboutFeatures(ctx context.Context) ([]model.AboutContent, error) {
	items, err := s.repo.GetActiveByType(ctx, "about_feature")
	if err != nil {
		return nil, err
	}
	return mapToAbout(items), nil
}

func (s *ContentService) GetTeamMembers(ctx context.Context) ([]repository.Content, error) {
	return s.repo.GetActiveByType(ctx, "team_member")
}

func (s *ContentService) GetTestimonials(ctx context.Context) ([]repository.Content, error) {
	return s.repo.GetActiveByType(ctx, "about_testimonial")
}

func (s *ContentService) GetStats(ctx context.Context) ([]repository.Content, error) {
	return s.repo.GetActiveByType(ctx, "about_stat")
}

func mapToResources(items []repository.Content) []model.Resource {
	out := make([]model.Resource, 0, len(items))

	for _, c := range items {
		out = append(out, model.Resource{
			BaseContent: mapBaseContent(c),
			Description: derefString(c.Description),
			FileURL:     derefString(c.ImageURL),
		})
	}

	return out
}

func (s *ContentService) GetResources(ctx context.Context) ([]model.Resource, error) {
	items, err := s.repo.GetActiveByType(ctx, "resource")
	if err != nil {
		return nil, err
	}
	return mapToResources(items), nil
}

func (s *ContentService) GetTopResources(ctx context.Context) ([]model.Resource, error) {
	items, err := s.repo.GetTopResources(ctx, 6)
	if err != nil {
		return nil, err
	}
	return mapToResources(items), nil
}

func mapToEvents(items []repository.EventContent) ([]model.Event, error) {
	out := make([]model.Event, 0, len(items))

	for _, e := range items {
		var eventDate time.Time
		if e.EventDate != "" {
			parsed, err := time.Parse("2006-01-02", e.EventDate)
			if err == nil {
				eventDate = parsed
			}
			// If parsing fails, eventDate stays as zero time (will be handled in frontend)
		}

		out = append(out, model.Event{
			BaseContent: model.BaseContent{
				ID:    e.ID.String(),
				Title: e.Title,
				Order: e.OrderIndex,
			},
			Description:      derefString(e.Description),
			EventDate:        eventDate,
			ImageURL:         derefString(e.ImageURL),
			Status:           "upcoming",
			RegistrationLink: e.RegistrationLink,
			Venue:            e.Venue,
			Price:            e.Price,
		})
	}

	return out, nil
}

func (s *ContentService) GetUpcomingEvents(ctx context.Context) ([]model.Event, error) {
	items, err := s.repo.GetUpcomingEvents(ctx, 3)
	if err != nil {
		return nil, err
	}
	return mapToEvents(items)
}

func (s *ContentService) GetAllEvents(ctx context.Context) ([]model.Event, error) {
	items, err := s.repo.GetAllEvents(ctx)
	if err != nil {
		return nil, err
	}
	return mapToEvents(items)
}

func (s *ContentService) GetAllContent(ctx context.Context) ([]repository.Content, error) {
	return s.repo.GetAll(ctx)
}

func (s *ContentService) CreateContent(ctx context.Context, c *repository.Content) error {
	return s.repo.Create(ctx, c)
}

func (s *ContentService) UpdateContent(ctx context.Context, id uuid.UUID, c *repository.Content) error {
	return s.repo.Update(ctx, id, c)
}

func (s *ContentService) DeleteContent(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}
