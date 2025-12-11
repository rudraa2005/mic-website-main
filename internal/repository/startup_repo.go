package repository

import(
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rudraa2005/mic-website-main/backend/internal/model"
)

type StartupRepository struct{
	db *pgxpool.Pool
}

func NewStartupRepository(db *pgxpool.Pool) *StartupRepository{
	return &StartupRepository{db}
}

func (r *StartupRepository) Create(ctx context.Context, s *model.Startup) error {
	query:=`
		INSERT INTO startup_ideas (owner_id, title, description, department)
        VALUES ($1, $2, $3, $4)
        RETURNING id, stage, created_at, updated_at;
		`
		return r.db.QueryRow(ctx,query,s.OwnerID,s.Title,s.Description,s.Department).Scan(&s.ID, &s.Stage, &s.CreatedAt, &s.UpdatedAt)
}

func (r *StartupRepository) GetByID(ctx context.Context, id string) (*model.Startup, error) {
    s := model.Startup{}
    query := `SELECT id, owner_id, title, description, stage, department, created_at, updated_at 
              FROM startup_ideas WHERE id=$1`

    err := r.db.QueryRow(ctx, query, id).
        Scan(&s.ID, &s.OwnerID, &s.Title, &s.Description, &s.Stage, &s.Department, &s.CreatedAt, &s.UpdatedAt)

    if err != nil {
        return nil, err
    }
    return &s, nil
}

func (r *StartupRepository) ListByOwner(ctx context.Context, ownerID string) ([]model.Startup, error) {
    query := `SELECT id, owner_id, title, description, stage, department, created_at, updated_at 
              FROM startup_ideas WHERE owner_id=$1 ORDER BY created_at DESC`

    rows, err := r.db.Query(ctx, query, ownerID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    startups := []model.Startup{}
    for rows.Next() {
        var s model.Startup
        rows.Scan(&s.ID, &s.OwnerID, &s.Title, &s.Description, &s.Stage, &s.Department, &s.CreatedAt, &s.UpdatedAt)
        startups = append(startups, s)
    }
    return startups, nil
}