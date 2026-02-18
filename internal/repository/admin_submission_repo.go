package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AdminSubmission struct {
	ID              string    `json:"id"`
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	Student         string    `json:"student"`
	FilePath        *string   `json:"file_path"`
	CreatedAt       time.Time `json:"submitted_on"`
	Status          string    `json:"status"`
	Tags            []string  `json:"tags"`
	Domain          *string   `json:"domain"`
	AssignedFaculty []string  `json:"assigned_faculty"`
}

type FacultyAssignment struct {
	ID           string    `json:"id"`
	SubmissionID string    `json:"submission_id"`
	FacultyID    string    `json:"faculty_id"`
	FacultyName  string    `json:"faculty_name"`
	AssignedAt   time.Time `json:"assigned_at"`
}

type AdminSubmissionRepo struct {
	db *pgxpool.Pool
}

func NewAdminSubmissionRepo(db *pgxpool.Pool) *AdminSubmissionRepo {
	return &AdminSubmissionRepo{db: db}
}

// GetPendingSubmissions returns submissions that students have submitted (status='submitted')
// These are waiting for admin review before going to faculty
func (r *AdminSubmissionRepo) GetPendingSubmissions(ctx context.Context) ([]AdminSubmission, error) {
	query := `
		SELECT
			s.submission_id,
			s.title,
			s.description,
			u.name,
			s.file_path,
			s.created_at,
			s.status
		FROM submissions s
		JOIN users u ON s.user_id = u.id
		WHERE s.status = 'submitted'
		ORDER BY s.created_at DESC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []AdminSubmission

	for rows.Next() {
		var s AdminSubmission
		err := rows.Scan(
			&s.ID,
			&s.Title,
			&s.Description,
			&s.Student,
			&s.FilePath,
			&s.CreatedAt,
			&s.Status,
		)
		if err != nil {
			return nil, err
		}
		res = append(res, s)
	}

	return res, nil
}

// ApproveForFaculty moves the submission to 'admin_approved' status
// so faculty can now review it
func (r *AdminSubmissionRepo) ApproveForFaculty(ctx context.Context, submissionID string) error {
	cmd, err := r.db.Exec(ctx, `
		UPDATE submissions
		SET status = 'admin_approved',
		    updated_at = now()
		WHERE submission_id = $1
		  AND status = 'submitted'
	`, submissionID)

	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return errors.New("submission not found or already processed")
	}

	return nil
}

// RejectSubmission marks submission as rejected by admin
func (r *AdminSubmissionRepo) RejectSubmission(ctx context.Context, submissionID string, reason string) error {
	cmd, err := r.db.Exec(ctx, `
		UPDATE submissions
		SET status = 'admin_rejected',
		    updated_at = now()
		WHERE submission_id = $1
		  AND status = 'submitted'
	`, submissionID)

	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return errors.New("submission not found or already processed")
	}

	return nil
}

// GetAllSubmissions returns ALL submissions for admin view (not just pending)
func (r *AdminSubmissionRepo) GetAllSubmissions(ctx context.Context) ([]AdminSubmission, error) {
	query := `
		SELECT
			s.submission_id,
			s.title,
			s.description,
			u.name,
			s.file_path,
			s.created_at,
			s.status,
			COALESCE(s.tags, '{}'),
			s.domain
		FROM submissions s
		JOIN users u ON s.user_id = u.id
		WHERE s.status != 'draft'
		ORDER BY s.created_at DESC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []AdminSubmission

	for rows.Next() {
		var s AdminSubmission
		err := rows.Scan(
			&s.ID,
			&s.Title,
			&s.Description,
			&s.Student,
			&s.FilePath,
			&s.CreatedAt,
			&s.Status,
			&s.Tags,
			&s.Domain,
		)
		if err != nil {
			return nil, err
		}
		res = append(res, s)
	}

	return res, nil
}

// AssignFacultyToSubmission assigns a faculty member to review a submission
func (r *AdminSubmissionRepo) AssignFacultyToSubmission(ctx context.Context, submissionID, facultyID, assignedByID string) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO submission_faculty (submission_id, faculty_id, assigned_by)
		VALUES ($1, $2, $3)
		ON CONFLICT (submission_id, faculty_id) DO NOTHING
	`, submissionID, facultyID, assignedByID)

	return err
}

// RemoveFacultyFromSubmission removes a faculty assignment
func (r *AdminSubmissionRepo) RemoveFacultyFromSubmission(ctx context.Context, submissionID, facultyID string) error {
	_, err := r.db.Exec(ctx, `
		DELETE FROM submission_faculty
		WHERE submission_id = $1 AND faculty_id = $2
	`, submissionID, facultyID)

	return err
}

// GetAssignedFaculty returns all faculty assigned to a submission
func (r *AdminSubmissionRepo) GetAssignedFaculty(ctx context.Context, submissionID string) ([]FacultyAssignment, error) {
	rows, err := r.db.Query(ctx, `
		SELECT sf.id, sf.submission_id, sf.faculty_id, u.name, sf.assigned_at
		FROM submission_faculty sf
		JOIN users u ON sf.faculty_id = u.id
		WHERE sf.submission_id = $1
		ORDER BY sf.assigned_at DESC
	`, submissionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []FacultyAssignment
	for rows.Next() {
		var fa FacultyAssignment
		if err := rows.Scan(&fa.ID, &fa.SubmissionID, &fa.FacultyID, &fa.FacultyName, &fa.AssignedAt); err != nil {
			return nil, err
		}
		res = append(res, fa)
	}
	return res, nil
}

// UpdateSubmissionTags updates the tags and domain for a submission
func (r *AdminSubmissionRepo) UpdateSubmissionTags(ctx context.Context, submissionID string, tags []string, domain string) error {
	_, err := r.db.Exec(ctx, `
		UPDATE submissions
		SET tags = $2, domain = $3, updated_at = now()
		WHERE submission_id = $1
	`, submissionID, tags, domain)

	return err
}
