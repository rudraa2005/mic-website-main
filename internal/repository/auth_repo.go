package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rudraa2005/mic-website-main/backend/internal/model"
)

type AuthRepository struct {
	db *pgxpool.Pool
}

func NewAuthRepository(db *pgxpool.Pool) *AuthRepository {
	return &AuthRepository{
		db: db,
	}
}

func (r *AuthRepository) Create(ctx context.Context, user *model.User) (string, error) {
	query := `
		INSERT INTO users (email, password_hash, role, name)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	err := r.db.QueryRow(
		context.Background(),
		query,
		user.Email,
		user.HashedPassword,
		user.Role,
		user.Name,
	).Scan(&user.UserID)

	return user.UserID, err
}

func (r *AuthRepository) FindByEmail(email string) (*model.User, error) {
	query := `
		SELECT id, email, password_hash, role, name
		FROM users
		WHERE email = $1
	`

	var user model.User
	err := r.db.QueryRow(context.Background(), query, email).
		Scan(&user.UserID, &user.Email, &user.HashedPassword, &user.Role, &user.Name)

	if err != nil {
		return nil, errors.New("user not found")
	}

	return &user, nil
}

func (r *AuthRepository) UpdatePassword(ctx context.Context, userID string, hashedPassword string) error {
	query := `
	UPDATE users 
	SET password_hash = $2, updated_at = now()
	WHERE id = $1
	`
	cmdTag, err := r.db.Exec(ctx, query, userID, hashedPassword)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return errors.New("user not found")
	}

	return nil
}

func (r *AuthRepository) GetByID(
	ctx context.Context,
	userID string,
) (*model.User, error) {

	query := `
		SELECT
			id,
			email,
			password_hash,
			role
		FROM users
		WHERE id = $1
	`

	var user model.User
	err := r.db.QueryRow(ctx, query, userID).
		Scan(
			&user.UserID,
			&user.Email,
			&user.HashedPassword,
			&user.Role,
		)

	if err != nil {
		return nil, errors.New("user not found")
	}

	return &user, nil
}
