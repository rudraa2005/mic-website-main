package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type FacultyProgress struct {
	SubmissionID string    `json:"submission_id"`
	Title        string    `json:"title"`
	Student      string    `json:"student"`
	AcceptedAt   time.Time `json:"accepted_at"`
	Stage        string    `json:"stage"`
	Progress     int       `json:"progress_percent"`
	Domain       string    `json:"domain"`
}

type FacultyProgressRepository struct {
	db *pgxpool.Pool
}

func NewFacultyProgressRepository(db *pgxpool.Pool) *FacultyProgressRepository {
	return &FacultyProgressRepository{db: db}
}

func (r *FacultyProgressRepository) GetByFaculty(
	ctx context.Context,
	facultyID string,
) ([]FacultyProgress, error) {

	// Since we no longer have faculty_id in work table, we join via submissions
	query := `
		SELECT
			s.submission_id,
			s.title,
			u.name,
			w.updated_at,
			w.stage,
			w.progress_percent
		FROM work w
		JOIN submissions s ON s.submission_id = w.submission_id
		JOIN users u ON u.id = s.user_id
		WHERE s.status = 'approved'
		ORDER BY w.updated_at DESC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []FacultyProgress

	for rows.Next() {
		var f FacultyProgress

		if err := rows.Scan(
			&f.SubmissionID,
			&f.Title,
			&f.Student,
			&f.AcceptedAt,
			&f.Stage,
			&f.Progress,
		); err != nil {
			return nil, err
		}
		f.Domain = "Unspecified"

		res = append(res, f)
	}

	return res, nil
}

func (r *FacultyProgressRepository) GetBySubmission(
	ctx context.Context,
	facultyID string,
	submissionID string,
) (*FacultyProgress, error) {

	query := `
		SELECT
			s.submission_id,
			s.title,
			u.name,
			w.updated_at,
			w.stage,
			w.progress_percent
		FROM work w
		JOIN submissions s ON s.submission_id = w.submission_id
		JOIN users u ON u.id = s.user_id
		WHERE s.submission_id = $1
		  AND s.status = 'approved'
		LIMIT 1
	`

	var f FacultyProgress
	err := r.db.QueryRow(ctx, query, submissionID).Scan(
		&f.SubmissionID,
		&f.Title,
		&f.Student,
		&f.AcceptedAt,
		&f.Stage,
		&f.Progress,
	)

	if err != nil {
		return nil, err
	}

	f.Domain = "Unspecified"
	return &f, nil
}
func (r *FacultyProgressRepository) UpdateProgress(ctx context.Context, submissionID string, stage string, progress int) error {
	query := `
		UPDATE work
		SET stage = $1, progress_percent = $2, updated_at = NOW()
		WHERE submission_id = $3
	`
	_, err := r.db.Exec(ctx, query, stage, progress, submissionID)
	return err
}

func (r *FacultyProgressRepository) LinkCompany(ctx context.Context, submissionID string, companyID string) error {
	query := `
		UPDATE work
		SET company_id = $1, updated_at = NOW()
		WHERE submission_id = $2
	`
	var cid interface{} = companyID
	if companyID == "" {
		cid = nil
	}
	_, err := r.db.Exec(ctx, query, cid, submissionID)
	return err
}
