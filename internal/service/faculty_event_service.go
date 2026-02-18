package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/rudraa2005/mic-website-main/backend/internal/repository"
)

type FacultyEvent struct {
	InvitationID string    `json:"id"`
	EventDate    string    `json:"event_date"`
	Title        string    `json:"title"`
	Date         string    `json:"date"`
	Venue        string    `json:"location"`
	Price        string    `json:"price"`
	InvitedAt    time.Time `json:"invited_at"`

	Status string `json:"status"` // pending | accepted | declined
}

var (
	ErrInvalidRSVPStatus = errors.New("invalid rsvp status")
	ErrUnauthorizedRSVP  = errors.New("unauthorized invitation access")
)

type EventInvitationService struct {
	repo *repository.EventInvitationRepository
}

func NewEventInvitationService(
	repo *repository.EventInvitationRepository,
) *EventInvitationService {
	return &EventInvitationService{repo: repo}
}

func (s *EventInvitationService) GetFacultyEvents(
	ctx context.Context,
	facultyID string,
) ([]FacultyEvent, error) {

	fid, err := uuid.Parse(facultyID)
	if err != nil {
		return nil, err
	}

	rows, err := s.repo.GetByFacultyID(ctx, fid)
	if err != nil {
		return nil, err
	}

	events := make([]FacultyEvent, 0, len(rows))

	for _, r := range rows {
		events = append(events, FacultyEvent{
			InvitationID: r.InvitationID.String(),
			Title:        r.Title,
			EventDate:    r.EventDate,
			Venue:        r.Venue,
			Price:        r.Price,
			InvitedAt:    r.InvitedAt,
			Status:       r.Status,
		})
	}

	return events, nil
}

func (s *EventInvitationService) UpdateRSVP(
	ctx context.Context,
	invitationID string,
	facultyID string,
	status string,
) error {

	if status != "accepted" && status != "declined" {
		return ErrInvalidRSVPStatus
	}

	iid, err := uuid.Parse(invitationID)
	fid, err := uuid.Parse(facultyID)
	if err != nil {
		return err
	}
	return s.repo.UpdateStatus(ctx, iid, fid, status)
}

func CalculateEventStats(events []FacultyEvent) (accepted, pending, declined int) {
	for _, e := range events {
		switch e.Status {
		case "accepted":
			accepted++
		case "pending":
			pending++
		case "declined":
			declined++
		}
	}
	return
}
