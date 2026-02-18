package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rudraa2005/mic-website-main/backend/internal/model"
)

type SettingsRepo struct {
	db *pgxpool.Pool
}

func NewSettingsRepo(db *pgxpool.Pool) *SettingsRepo {
	return &SettingsRepo{
		db: db,
	}
}
func (r *SettingsRepo) GetByUserID(
	ctx context.Context,
	userID string,
) (*model.Settings, error) {

	query := `
		SELECT
			user_id,
			theme,
			email_notifications,
			feedback_alerts,
			application_updates,
			newsletter,
			created_at,
			updated_at
		FROM settings
		WHERE user_id = $1
	`

	var s model.Settings
	err := r.db.QueryRow(ctx, query, userID).
		Scan(
			&s.UserID,
			&s.Theme,
			&s.EmailNotifications,
			&s.FeedbackAlerts,
			&s.ApplicationUpdates,
			&s.Newsletter,
			&s.CreatedAt,
			&s.UpdatedAt,
		)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, pgx.ErrNoRows
		}
		return nil, err
	}

	return &s, nil
}

func (r *SettingsRepo) Update(
	ctx context.Context,
	s model.Settings,
) error {

	query := `
		INSERT INTO settings (
			user_id,
			theme,
			email_notifications,
			feedback_alerts,
			application_updates,
			newsletter
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (user_id) DO UPDATE
		SET
			theme = EXCLUDED.theme,
			email_notifications = EXCLUDED.email_notifications,
			feedback_alerts = EXCLUDED.feedback_alerts,
			application_updates = EXCLUDED.application_updates,
			newsletter = EXCLUDED.newsletter,
			updated_at = now()
	`

	_, err := r.db.Exec(
		ctx,
		query,
		s.UserID,
		s.Theme,
		s.EmailNotifications,
		s.FeedbackAlerts,
		s.ApplicationUpdates,
		s.Newsletter,
	)

	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			return pgErr
		}
		return err
	}

	return nil
}

func (r *SettingsRepo) CreateDefaults(ctx context.Context, userID string) error {
	query := `
		INSERT INTO settings (user_id)
		VALUES ($1)
		ON CONFLICT (user_id) DO NOTHING
	`
	_, err := r.db.Exec(ctx, query, userID)
	return err
}
