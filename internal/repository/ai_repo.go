package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AIRepo struct {
	db *pgxpool.Pool
}

func NewAIRepo(db *pgxpool.Pool) *AIRepo {
	return &AIRepo{db: db}
}

func (r *AIRepo) CreateDraft(
	ctx context.Context,
	submissionID string,
	userID string,
	filePath string,
) error {
	_, err := r.db.Exec(ctx, `
        INSERT INTO submission_ai_insights (submission_id, user_id, file_path, status)
        VALUES ($1, $2, $3, 'processing')
        ON CONFLICT (submission_id) DO NOTHING
    `, submissionID, userID, filePath)

	return err
}
func (r *AIRepo) SaveInsights(
	ctx context.Context,
	submissionID string,
	insights string,
) error {
	_, err := r.db.Exec(ctx, `
        UPDATE submission_ai_insights
        SET insights = $1, status = 'completed'
        WHERE submission_id = $2
    `, insights, submissionID)

	return err
}

func (r *AIRepo) GetInsights(ctx context.Context, submissionID string) (string, string, error) {
	var insights, status string
	err := r.db.QueryRow(ctx, `
		SELECT insights, status
		FROM submission_ai_insights
		WHERE submission_id = $1
	`, submissionID).Scan(&insights, &status)
	return insights, status, err
}
