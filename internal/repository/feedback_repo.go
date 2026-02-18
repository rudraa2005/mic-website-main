package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rudraa2005/mic-website-main/backend/internal/model"
)

type FeedbackRepo struct {
	db *pgxpool.Pool
}

func NewFeedbackRepo(db *pgxpool.Pool) *FeedbackRepo {
	return &FeedbackRepo{db: db}
}

func (r *FeedbackRepo) GetByUserID(
	ctx context.Context,
	userID string,
) ([]model.Feedback, error) {

	query := `
		SELECT
			f.feedback_id,
			f.submission_id,
			f.faculty_id,
			f.faculty_name,
			f.faculty_title,
			f.faculty_field,
			f.overall_feedback,
			f.strengths,
			f.recommendations,
			f.rating,
			f.status,
			f.created_at,
			f.updated_at
		FROM feedbacks f
		JOIN submissions s ON s.submission_id = f.submission_id
		WHERE s.user_id = $1
		ORDER BY f.updated_at DESC
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var feedbacks []model.Feedback

	for rows.Next() {
		var f model.Feedback
		err := rows.Scan(
			&f.FeedbackID,
			&f.SubmissionID,
			&f.FacultyID,
			&f.FacultyName,
			&f.FacultyTitle,
			&f.FacultyField,
			&f.OverallFeedback,
			&f.Strengths,
			&f.Recommendations,
			&f.Rating,
			&f.Status,
			&f.CreatedAt,
			&f.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		feedbacks = append(feedbacks, f)
	}

	return feedbacks, nil
}
func (r *FeedbackRepo) Create(ctx context.Context, f *model.Feedback) error {
	if f.FeedbackID == "" {
		f.FeedbackID = uuid.NewString()
	}
	if f.Strengths == nil {
		f.Strengths = []string{}
	}
	if f.Recommendations == nil {
		f.Recommendations = []string{}
	}

	query := `
		INSERT INTO feedbacks (
			feedback_id,
			submission_id,
			faculty_id,
			faculty_name,
			faculty_title,
			faculty_field,
			overall_feedback,
			strengths,
			recommendations,
			rating,
			status,
			created_at,
			updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, NOW(), NOW())
		RETURNING feedback_id
	`

	return r.db.QueryRow(ctx, query,
		f.FeedbackID,
		f.SubmissionID,
		f.FacultyID,
		f.FacultyName,
		f.FacultyTitle,
		f.FacultyField,
		f.OverallFeedback,
		f.Strengths,
		f.Recommendations,
		f.Rating,
		f.Status,
	).Scan(&f.FeedbackID)
}
