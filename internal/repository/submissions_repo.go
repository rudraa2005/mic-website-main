package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rudraa2005/mic-website-main/backend/internal/model"
)

type SubmissionsRepo struct {
	db *pgxpool.Pool
}

func NewSubmissionsRepo(db *pgxpool.Pool) *SubmissionsRepo {
	return &SubmissionsRepo{
		db: db,
	}
}
func (r *SubmissionsRepo) GetBySubmissionID(
	ctx context.Context,
	submissionID string,
) (*model.Submission, error) {

	query := `
		SELECT 
			title,
			description,
			file_path,
			status,
			submission_id,
			stage,
			user_id,
			created_at,
			updated_at
		FROM submissions
		WHERE submission_id = $1
	`
	row, err := r.db.Query(ctx, query, submissionID)
	if err != nil {
		return nil, err
	}
	defer row.Close()

	var s model.Submission
	if row.Next() {
		err := row.Scan(
			&s.Title,
			&s.Description,
			&s.FilePath,
			&s.Status,
			&s.SubmissionID,
			&s.Stage,
			&s.UserID,
			&s.CreatedAt,
			&s.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		return &s, nil
	}
	return nil, errors.New("Submission Not Found!")
}

func (r *SubmissionsRepo) GetByUserID(
	ctx context.Context,
	userID string,
) ([]model.Submission, error) {

	query := `
		SELECT
			submission_id,
			user_id,
			title,
			description,
			file_path,
			status,
			stage,
			created_at,
			updated_at
		FROM submissions
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var submissions []model.Submission

	for rows.Next() {
		var s model.Submission
		err := rows.Scan(
			&s.SubmissionID,
			&s.UserID,
			&s.Title,
			&s.Description,
			&s.FilePath,
			&s.Status,
			&s.Stage,
			&s.CreatedAt,
			&s.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		submissions = append(submissions, s)
	}

	return submissions, nil
}

func (r *SubmissionsRepo) Create(
	ctx context.Context,
	s *model.Submission,
) error {

	query := `
		INSERT INTO submissions (
			submission_id,
			user_id,
			title,
			description,
			file_path,
			status,
			stage
		)
		VALUES ($1, $2, $3, $4, $5, 'draft',$6)
	`

	_, err := r.db.Exec(
		ctx,
		query,
		s.SubmissionID,
		s.UserID,
		s.Title,
		s.Description,
		s.FilePath,
		s.Stage,
	)

	return err
}

func (r *SubmissionsRepo) UpdateDraft(
	ctx context.Context,
	s *model.Submission,

) error {

	query := `
		UPDATE submissions
		SET
			title = $1,
			description = $2,
			file_path = COALESCE($3, file_path),
			updated_at = now()
		WHERE user_id = $4
		  AND submission_id = $5
		  AND status = 'draft'
	`

	cmdTag, err := r.db.Exec(
		ctx,
		query,
		s.Title,
		s.Description,
		s.FilePath,
		s.UserID,
		s.SubmissionID,
	)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return errors.New("startup cannot be updated (not owner or not draft)")
	}

	return nil
}

func (r *SubmissionsRepo) Delete(
	ctx context.Context,
	submissionID string,
	userID string,
) error {

	query := `
		DELETE FROM submissions
		WHERE user_id = $1
		  AND submission_id = $2
	`

	cmd, err := r.db.Exec(ctx, query, userID, submissionID)
	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return errors.New("submission not found or not owned by user")
	}

	return nil
}

func (r *SubmissionsRepo) MarkSubmitted(
	ctx context.Context,
	submissionID string,
	userID string,
) error {

	query := `
        UPDATE submissions
        SET status = 'submitted',
            updated_at = now()
        WHERE submission_id = $1
          AND user_id = $2
          AND status = 'draft'
    `

	cmd, err := r.db.Exec(ctx, query, submissionID, userID)
	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return errors.New("submission cannot be submitted")
	}

	return nil
}

func (r *SubmissionsRepo) AttachFile(
	ctx context.Context,
	submissionID string,
	userID string,
	filePath string,
) error {

	query := `
		UPDATE submissions
		SET file_path = $1,
		    updated_at = now()
		WHERE submission_id = $2
			AND user_id = $3
	`

	_, err := r.db.Exec(ctx, query, filePath, submissionID, userID)
	return err
}

func (r *SubmissionsRepo) UpdateStatus(
	ctx context.Context,
	submissionID string,
	newStatus string,
	newStage string,
) error {

	query := `
		UPDATE submissions
		SET status = $1,
		    stage = $2,
		    updated_at = now()
		WHERE submission_id = $3
	`

	cmd, err := r.db.Exec(ctx, query, newStatus, newStage, submissionID)
	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return errors.New("submission not found")
	}

	return nil
}

func (r *SubmissionsRepo) GetIncubationPipeline(ctx context.Context) ([]model.Submission, error) {
	query := `
		SELECT 
			w.title,
			w.description,
			s.file_path,
			s.status,
			s.submission_id,
			w.stage,
			s.user_id,
			w.created_at,
			w.updated_at,
			c.name as company_name,
			c.logo_url as company_logo
		FROM work w
		JOIN submissions s ON w.submission_id = s.submission_id
		LEFT JOIN companies c ON w.company_id = c.id
		ORDER BY w.created_at DESC
	`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var submissions []model.Submission
	for rows.Next() {
		var s model.Submission
		err := rows.Scan(
			&s.Title,
			&s.Description,
			&s.FilePath,
			&s.Status,
			&s.SubmissionID,
			&s.Stage,
			&s.UserID,
			&s.CreatedAt,
			&s.UpdatedAt,
			&s.CompanyName,
			&s.CompanyLogo,
		)
		if err != nil {
			return nil, err
		}
		submissions = append(submissions, s)
	}
	return submissions, nil
}
