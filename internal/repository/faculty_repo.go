package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rudraa2005/mic-website-main/backend/internal/model"
)

type FacultyRepository struct {
	db *pgxpool.Pool
}

func NewFacultyRepository(db *pgxpool.Pool) *FacultyRepository {
	return &FacultyRepository{db: db}
}

func (r *FacultyRepository) FindByEmail(email string) (*model.Faculty, error) {
	query := `
		SELECT id, name, email, password_hash, role
		FROM faculty
		WHERE email = $1
	`

	var f model.Faculty

	err := r.db.QueryRow(context.Background(), query, email).Scan(
		&f.ID,
		&f.Name,
		&f.Email,
		&f.PasswordHash,
		&f.Role,
	)

	if err != nil {
		return nil, err
	}

	return &f, nil
}
