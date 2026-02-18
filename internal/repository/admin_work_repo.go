package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type WorkItem struct {
	ID              string    `json:"id"`
	SubmissionID    string    `json:"submission_id"`
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	Stage           string    `json:"stage"`
	ProgressPercent int       `json:"progress_percent"`
	CompanyID       *string   `json:"company_id"`
	CompanyName     *string   `json:"company_name"`
	CompanyLogo     *string   `json:"company_logo"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type Company struct {
	ID      string  `json:"id"`
	Name    string  `json:"name"`
	LogoURL *string `json:"logo_url"`
}

type AdminWorkRepo struct {
	db *pgxpool.Pool
}

func NewAdminWorkRepo(db *pgxpool.Pool) *AdminWorkRepo {
	return &AdminWorkRepo{db: db}
}

// GetAllWork returns all work items with company info
func (r *AdminWorkRepo) GetAllWork(ctx context.Context) ([]WorkItem, error) {
	query := `
		SELECT
			w.id,
			w.submission_id,
			w.title,
			w.description,
			w.stage,
			w.progress_percent,
			w.company_id,
			c.name,
			c.logo_url,
			w.created_at,
			w.updated_at
		FROM work w
		LEFT JOIN companies c ON w.company_id = c.id
		ORDER BY w.created_at DESC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []WorkItem
	for rows.Next() {
		var w WorkItem
		err := rows.Scan(
			&w.ID,
			&w.SubmissionID,
			&w.Title,
			&w.Description,
			&w.Stage,
			&w.ProgressPercent,
			&w.CompanyID,
			&w.CompanyName,
			&w.CompanyLogo,
			&w.CreatedAt,
			&w.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		res = append(res, w)
	}

	return res, nil
}

// UpdateWork updates a work item
func (r *AdminWorkRepo) UpdateWork(ctx context.Context, workID, title, description, stage string, progressPercent int, companyID *string) error {
	_, err := r.db.Exec(ctx, `
		UPDATE work
		SET title = $2, description = $3, stage = $4, progress_percent = $5, company_id = $6, updated_at = now()
		WHERE id = $1
	`, workID, title, description, stage, progressPercent, companyID)

	return err
}

// DeleteWork removes a work item
func (r *AdminWorkRepo) DeleteWork(ctx context.Context, workID string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM work WHERE id = $1`, workID)
	return err
}

// GetCompanies returns all companies
func (r *AdminWorkRepo) GetCompanies(ctx context.Context) ([]Company, error) {
	rows, err := r.db.Query(ctx, `SELECT id, name, logo_url FROM companies ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []Company
	for rows.Next() {
		var c Company
		if err := rows.Scan(&c.ID, &c.Name, &c.LogoURL); err != nil {
			return nil, err
		}
		res = append(res, c)
	}

	return res, nil
}

// AddCompany creates a new company
func (r *AdminWorkRepo) AddCompany(ctx context.Context, name string, logoURL *string) (string, error) {
	var id string
	err := r.db.QueryRow(ctx, `
		INSERT INTO companies (name, logo_url)
		VALUES ($1, $2)
		RETURNING id
	`, name, logoURL).Scan(&id)

	return id, err
}

// DeleteCompany removes a company
func (r *AdminWorkRepo) DeleteCompany(ctx context.Context, companyID string) error {
	// First clear company_id from any work items
	_, _ = r.db.Exec(ctx, `UPDATE work SET company_id = NULL WHERE company_id = $1`, companyID)

	// Then delete the company
	_, err := r.db.Exec(ctx, `DELETE FROM companies WHERE id = $1`, companyID)
	return err
}
