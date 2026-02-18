package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rudraa2005/mic-website-main/backend/internal/model"
)

type CompanyRepo struct {
	db *pgxpool.Pool
}

func NewCompanyRepo(db *pgxpool.Pool) *CompanyRepo {
	return &CompanyRepo{db: db}
}

func (r *CompanyRepo) GetAll(ctx context.Context) ([]model.Company, error) {
	query := `SELECT id, name, logo_url, created_at FROM companies ORDER BY name ASC`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var companies []model.Company
	for rows.Next() {
		var c model.Company
		err := rows.Scan(&c.ID, &c.Name, &c.LogoURL, &c.CreatedAt)
		if err != nil {
			return nil, err
		}
		companies = append(companies, c)
	}
	return companies, nil
}

func (r *CompanyRepo) GetByID(ctx context.Context, id string) (*model.Company, error) {
	var c model.Company
	query := `SELECT id, name, logo_url, created_at FROM companies WHERE id = $1`
	err := r.db.QueryRow(ctx, query, id).Scan(&c.ID, &c.Name, &c.LogoURL, &c.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &c, nil
}
