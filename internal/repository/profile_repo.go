package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rudraa2005/mic-website-main/backend/internal/model"
)

type ProfileRepository struct {
	db *pgxpool.Pool
}

func NewProfileRepo(db *pgxpool.Pool) *ProfileRepository {
	return &ProfileRepository{
		db: db,
	}
}
func (r *ProfileRepository) GetUserByID(ctx context.Context, userID string) (*model.Profile, error) {
	query := `
	SELECT 
		COALESCE(p.user_id, u.id) as user_id,
		COALESCE(p.name, u.name) as name,
		COALESCE(p.email, u.email) as email,
		p.phone,
		p.photo_url,
		COALESCE(p.bio, '') as bio,
		COALESCE(p.created_at, u.created_at) as created_at,
		COALESCE(p.updated_at, u.updated_at) as updated_at
	FROM users u
	LEFT JOIN profiles p ON u.id = p.user_id
	WHERE u.id = $1
	`

	var p model.Profile
	err := r.db.QueryRow(ctx, query, userID).
		Scan(&p.UserID, &p.Name, &p.Email, &p.Phone, &p.PhotoURL, &p.Bio, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &p, nil
}

func (r *ProfileRepository) StoreUser(ctx context.Context, p model.Profile) error {
	query := `
		INSERT INTO profiles (user_id, name, phone, email, photo_url)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (user_id) DO UPDATE 
		SET name = EXCLUDED.name, email = EXCLUDED.email
	`
	_, err := r.db.Exec(ctx, query, p.UserID, p.Name, p.Phone, p.Email, p.PhotoURL)

	return err
}

func (r *ProfileRepository) UpdatePhotoURL(ctx context.Context, userID string, photoURL string) error {
	query := `
	UPDATE profiles
	SET photo_url = $2
	WHERE user_id = $1
	`

	_, err := r.db.Exec(ctx, query, userID, photoURL)

	return err
}
func (r *ProfileRepository) UpdateProfile(
	ctx context.Context,
	userID string,
	name string,
	phone *string,
	bio string,
) error {

	query := `
	INSERT INTO profiles (user_id, name, phone, bio)
	VALUES ($1, $2, $3, $4)
	ON CONFLICT (user_id) DO UPDATE SET
		name = EXCLUDED.name,
		phone = EXCLUDED.phone,
		bio = EXCLUDED.bio,
		updated_at = now()
	`
	_, err := r.db.Exec(
		ctx,
		query,
		userID,
		name,
		phone,
		bio,
	)

	return err
}
