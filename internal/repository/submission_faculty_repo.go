package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type FacultySubmission struct {
	ID              string    `json:"id"`
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	Student         string    `json:"student"`
	Email           string    `json:"email"`
	FilePath        *string   `json:"file_path"`
	CreatedAt       time.Time `json:"submitted_on"`
	Status          string    `json:"status"`
	Tags            []string  `json:"tags"`
	Domain          *string   `json:"domain"`
	Stage           *string   `json:"stage"`
	ProgressPercent *int      `json:"progress_percent"`
}

type FacultySubmissionRepo struct {
	db *pgxpool.Pool
}

func NewFacultySubmissionRepo(db *pgxpool.Pool) *FacultySubmissionRepo {
	return &FacultySubmissionRepo{db: db}
}

func (r *FacultySubmissionRepo) GetSubmitted(ctx context.Context) ([]FacultySubmission, error) {
	query := `
		SELECT
			s.submission_id,
			s.title,
			s.description,
			u.name,
			u.email,
			s.file_path,
			s.created_at,
			s.status,
			COALESCE(s.tags, '{}'),
			s.domain
		FROM submissions s
		JOIN users u ON s.user_id = u.id
		WHERE s.status IN ('admin_approved', 'approved', 'rejected')
		ORDER BY s.created_at DESC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []FacultySubmission

	for rows.Next() {
		var f FacultySubmission
		err := rows.Scan(
			&f.ID,
			&f.Title,
			&f.Description,
			&f.Student,
			&f.Email,
			&f.FilePath,
			&f.CreatedAt,
			&f.Status,
			&f.Tags,
			&f.Domain,
		)
		if err != nil {
			return nil, err
		}
		res = append(res, f)
	}

	return res, nil
}

func (r *FacultySubmissionRepo) GetByID(ctx context.Context, id string) (*FacultySubmission, error) {
	query := `
		SELECT
			s.submission_id,
			s.title,
			s.description,
			u.name,
			u.email,
			s.file_path,
			s.created_at,
			s.status,
			COALESCE(s.tags, '{}'),
			s.domain,
			w.stage,
			w.progress_percent
		FROM submissions s
		JOIN users u ON s.user_id = u.id
		LEFT JOIN work w ON w.submission_id = s.submission_id
		WHERE s.submission_id = $1
		  AND s.status IN ('admin_approved', 'approved', 'rejected')
	`

	var f FacultySubmission
	err := r.db.QueryRow(ctx, query, id).Scan(
		&f.ID,
		&f.Title,
		&f.Description,
		&f.Student,
		&f.Email,
		&f.FilePath,
		&f.CreatedAt,
		&f.Status,
		&f.Tags,
		&f.Domain,
		&f.Stage,
		&f.ProgressPercent,
	)
	if err != nil {
		return nil, err
	}

	return &f, nil
}

func (r *FacultySubmissionRepo) ApproveSubmission(
	ctx context.Context,
	submissionID string,
	facultyID string,
) (email string, title string, err error) {

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return "", "", err
	}
	defer tx.Rollback(ctx)

	// Fetch details for email
	var uEmail, sTitle, sDesc string
	err = tx.QueryRow(ctx, `
		SELECT u.email, s.title, s.description
		FROM submissions s
		JOIN users u ON s.user_id = u.id
		WHERE s.submission_id = $1
	`, submissionID).Scan(&uEmail, &sTitle, &sDesc)
	if err != nil {
		return "", "", err
	}

	// 1. Mark submission approved
	cmd, err := tx.Exec(ctx, `
		UPDATE submissions
		SET status = 'approved'
		WHERE submission_id = $1
		  AND status = 'admin_approved'
	`, submissionID)

	if err != nil {
		return "", "", err
	}

	if cmd.RowsAffected() == 0 {
		return "", "", errors.New("submission not found or already processed")
	}

	// 2. Insert into work table
	_, err = tx.Exec(ctx, `
		INSERT INTO work (
			submission_id,
			title,
			description,
			stage,
			progress_percent
		)
		VALUES ($1, $2, $3, 'under_incubation', 0)
		ON CONFLICT (submission_id) DO NOTHING
	`, submissionID, sTitle, sDesc)

	if err != nil {
		return "", "", err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return "", "", err
	}

	return uEmail, sTitle, nil
}

func (r *FacultySubmissionRepo) RejectSubmission(
	ctx context.Context,
	submissionID string,
	facultyID string,
) (email string, title string, err error) {

	// Fetch details for email
	var uEmail, sTitle string
	err = r.db.QueryRow(ctx, `
		SELECT u.email, s.title
		FROM submissions s
		JOIN users u ON s.user_id = u.id
		WHERE s.submission_id = $1
	`, submissionID).Scan(&uEmail, &sTitle)
	if err != nil {
		return "", "", err
	}

	cmd, err := r.db.Exec(ctx, `
		UPDATE submissions
		SET status = 'rejected'
		WHERE submission_id = $1
		  AND status = 'admin_approved'
	`, submissionID)

	if err != nil {
		return "", "", err
	}

	if cmd.RowsAffected() == 0 {
		return "", "", errors.New("submission not found or already processed")
	}

	return uEmail, sTitle, nil
}

// NeedsImprovement marks a submission as needing improvement
func (r *FacultySubmissionRepo) NeedsImprovement(
	ctx context.Context,
	submissionID string,
	facultyID string,
) (email string, title string, err error) {

	// Fetch details for email
	var uEmail, sTitle string
	err = r.db.QueryRow(ctx, `
		SELECT u.email, s.title
		FROM submissions s
		JOIN users u ON s.user_id = u.id
		WHERE s.submission_id = $1
	`, submissionID).Scan(&uEmail, &sTitle)
	if err != nil {
		return "", "", err
	}

	cmd, err := r.db.Exec(ctx, `
		UPDATE submissions
		SET status = 'needs_improvement'
		WHERE submission_id = $1
		  AND status = 'admin_approved'
	`, submissionID)

	if err != nil {
		return "", "", err
	}

	if cmd.RowsAffected() == 0 {
		return "", "", errors.New("submission not found or already processed")
	}

	return uEmail, sTitle, nil
}
