package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type FacultyEventInvitation struct {
	InvitationID uuid.UUID
	ContentID    uuid.UUID

	Title     string
	EventDate string
	Venue     string
	Price     string

	Status      string
	InvitedAt   time.Time
	RespondedAt *time.Time
}

type EventInvitationRepository struct {
	db *pgxpool.Pool
}

func NewEventInvitationRepository(db *pgxpool.Pool) *EventInvitationRepository {
	return &EventInvitationRepository{db: db}
}

func (r *EventInvitationRepository) GetByFacultyID(
	ctx context.Context,
	facultyID uuid.UUID,
) ([]FacultyEventInvitation, error) {

	query := `
		SELECT
			ei.id                         AS invitation_id,
			ei.status                     AS status,
			ei.invited_at                 AS invited_at,

			c.title                       AS title,
			c.content_data->>'event_date' AS event_date,
			c.content_data->>'venue'      AS venue,
			c.content_data->>'price'      AS price
		FROM event_invitations ei
		JOIN content c
		  ON c.id = ei.content_id
		WHERE
		  ei.faculty_id = $1
		  AND c.content_type = 'event'
		  AND c.is_active = true
		ORDER BY (c.content_data->>'event_date')::date ASC
	`

	rows, err := r.db.Query(ctx, query, facultyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []FacultyEventInvitation

	for rows.Next() {
		var r FacultyEventInvitation
		if err := rows.Scan(
			&r.InvitationID,
			&r.Status,
			&r.InvitedAt,
			&r.Title,
			&r.EventDate,
			&r.Venue,
			&r.Price,
		); err != nil {
			return nil, err
		}
		result = append(result, r)
	}

	return result, nil
}

func (r *EventInvitationRepository) UpdateStatus(
	ctx context.Context,
	invitationID uuid.UUID,
	facultyID uuid.UUID,
	status string,
) error {

	query := `
		UPDATE event_invitations
		SET status = $1,
		    responded_at = now()
		WHERE id = $2
		  AND faculty_id = $3
		  AND status = 'pending'
	`

	cmd, err := r.db.Exec(ctx, query, status, invitationID, facultyID)
	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return errors.New("invalid invitation or already responded")
	}

	return nil
}

func (r *EventInvitationRepository) Create(
	ctx context.Context,
	contentID uuid.UUID,
	facultyID uuid.UUID,
) error {

	query := `
		INSERT INTO event_invitations (content_id, faculty_id)
		VALUES ($1, $2)
		ON CONFLICT (content_id, faculty_id) DO NOTHING
	`

	_, err := r.db.Exec(ctx, query, contentID, facultyID)
	return err
}
