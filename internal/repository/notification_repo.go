package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rudraa2005/mic-website-main/backend/internal/model"
)

type notificationRepository struct {
	db *pgxpool.Pool
}

func NewNotificationRepository(db *pgxpool.Pool) *notificationRepository {
	return &notificationRepository{db: db}
}

func (r *notificationRepository) NotifyStatusChange(ctx context.Context, n *model.Notification) error {
	query := `
		INSERT INTO notifications (
			user_id, role, type, message, submission_id
		) VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.Exec(
		ctx,
		query,
		n.UserID,
		n.Role,
		n.Type,
		n.Body,
		n.SubmissionID,
	)
	return err
}

func (r *notificationRepository) GetByUser(ctx context.Context, userID string) ([]model.Notification, error) {
	query := `
		SELECT id, user_id, role, type, message, submission_id, is_read, created_at
		FROM notifications
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []model.Notification

	for rows.Next() {
		var n model.Notification
		err := rows.Scan(
			&n.ID,
			&n.UserID,
			&n.Role,
			&n.Type,
			&n.Body,
			&n.SubmissionID,
			&n.IsRead,
			&n.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, n)
	}

	return result, nil
}

func (r *notificationRepository) MarkAsRead(ctx context.Context, id string, userID string) error {
	query := `
		UPDATE notifications
		SET is_read = true
		WHERE id = $1 AND user_id = $2
	`
	_, err := r.db.Exec(ctx, query, id, userID)
	return err
}

func (r *notificationRepository) GetUnreadCount(ctx context.Context, userID string) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM notifications
		WHERE user_id = $1 AND is_read = false
	`
	var count int
	err := r.db.QueryRow(ctx, query, userID).Scan(&count)
	return count, err
}

func (r *notificationRepository) GetNotificationsByUser(
	ctx context.Context,
	userID string,
) ([]model.Notification, error) {
	return r.GetByUser(ctx, userID)
}

func (r *notificationRepository) MarkNotificationAsRead(
	ctx context.Context,
	id string,
	userID string,
) error {
	return r.MarkAsRead(ctx, id, userID)
}

func (r *notificationRepository) GetUnreadCountByUser(
	ctx context.Context,
	userID string,
) (int, error) {
	return r.GetUnreadCount(ctx, userID)
}
