package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AdminFacultyUser struct {
	ID    string
	Name  string
	Email string
}

type AdminFacultyRepository struct {
	db *pgxpool.Pool
}

func NewAdminFacultyRepository(db *pgxpool.Pool) *AdminFacultyRepository {
	return &AdminFacultyRepository{db: db}
}

// GetAllByRole fetches all users with a specific role
func (r *AdminFacultyRepository) GetAllByRole(ctx context.Context, role string) ([]AdminFacultyUser, error) {
	query := `
		SELECT id, name, email
		FROM users
		WHERE role = $1
		ORDER BY name ASC
	`

	rows, err := r.db.Query(ctx, query, role)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []AdminFacultyUser
	for rows.Next() {
		var u AdminFacultyUser
		if err := rows.Scan(&u.ID, &u.Name, &u.Email); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

// Create adds a new user with the specified role
func (r *AdminFacultyRepository) Create(ctx context.Context, name, email, hashedPassword, role string) error {
	query := `
		INSERT INTO users (name, email, password_hash, role)
		VALUES ($1, $2, $3, $4)
	`

	_, err := r.db.Exec(ctx, query, name, email, hashedPassword, role)
	return err
}

// Update modifies an existing user's name, email, and optionally password
func (r *AdminFacultyRepository) Update(ctx context.Context, id, name, email, hashedPassword string) error {
	var query string
	var err error

	if hashedPassword != "" {
		query = `
			UPDATE users
			SET name = $2, email = $3, password_hash = $4, updated_at = now()
			WHERE id = $1
		`
		_, err = r.db.Exec(ctx, query, id, name, email, hashedPassword)
	} else {
		query = `
			UPDATE users
			SET name = $2, email = $3, updated_at = now()
			WHERE id = $1
		`
		_, err = r.db.Exec(ctx, query, id, name, email)
	}

	return err
}

// Delete removes a user by ID
func (r *AdminFacultyRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM users WHERE id = $1`

	cmdTag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return errors.New("faculty not found")
	}

	return nil
}
