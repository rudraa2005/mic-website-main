package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rudraa2005/mic-website-main/backend/internal/model"
)

type QueryRepo struct {
	db *pgxpool.Pool
}

func NewQueryRepo(db *pgxpool.Pool) *QueryRepo {
	return &QueryRepo{db: db}
}

func (r *QueryRepo) Create(
	ctx context.Context,
	q *model.Query,
) error {

	query := `
		INSERT INTO queries (
			query_id,
			user_id,
			faculty_id,
			feedback_id,
			query,
			priority,
			status
		)
		VALUES ($1, $2, $3, $4, $5, $6, 'pending')
	`

	_, err := r.db.Exec(
		ctx,
		query,
		q.QueryID,
		q.UserID,
		q.FacultyID,
		q.FeedbackID,
		q.QueryText,
		q.Priority,
	)

	return err
}

func (r *QueryRepo) GetByUserID(
	ctx context.Context,
	userID string,
) ([]model.Query, error) {

	query := `
		SELECT
			query_id,
			user_id,
			faculty_id,
			feedback_id,
			query,
			priority,
			status,
			response,
			created_at,
			updated_at
		FROM queries
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var queries []model.Query

	for rows.Next() {
		var q model.Query
		err := rows.Scan(
			&q.QueryID,
			&q.UserID,
			&q.FacultyID,
			&q.FeedbackID,
			&q.QueryText,
			&q.Priority,
			&q.Status,
			&q.Response,
			&q.CreatedAt,
			&q.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		queries = append(queries, q)
	}

	return queries, nil
}
