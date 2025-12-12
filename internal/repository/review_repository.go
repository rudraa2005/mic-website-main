package repository

import (
    "context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Review struct {
    ID                string   `json:"id"`
    StartupID         string   `json:"startup_id"`

    ReviewerID        string   `json:"reviewer_id"`
    ReviewerName      string   `json:"reviewer_name"`
    ReviewerDesignation string `json:"reviewer_designation"`

    Rating            float32  `json:"rating"`
    Decision          string   `json:"decision"`
    Summary           string   `json:"summary"`

    Strengths         []string `json:"strengths"`
    Recommendations   []string `json:"recommendations"`

    CreatedAt         string   `json:"created_at"`
}

type ReviewRepository struct {
    pool *pgxpool.Pool
}

func NewReviewRepository(pool *pgxpool.Pool) *ReviewRepository {
    return &ReviewRepository{pool: pool}
}

func (r *ReviewRepository) CreateReview(ctx context.Context, rev Review) error {
    _, err := r.pool.Exec(ctx,
        `INSERT INTO reviews 
        (startup_id, reviewer_id, reviewer_name, reviewer_designation,
         rating, decision, summary, strengths, recommendations)
        VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`,
        rev.StartupID,
        rev.ReviewerID, rev.ReviewerName, rev.ReviewerDesignation,
        rev.Rating, rev.Decision, rev.Summary,
        rev.Strengths, rev.Recommendations,
    )
    return err
}

func (r *ReviewRepository) ListForStartup(ctx context.Context, startupID string) ([]Review, error) {
    rows, err := r.pool.Query(ctx,
        `SELECT id, startup_id,
                reviewer_id, reviewer_name, reviewer_designation,
                rating, decision, summary, strengths, recommendations,
                created_at
         FROM reviews 
         WHERE startup_id=$1 
         ORDER BY created_at DESC`,
        startupID,
    )
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var out []Review
    for rows.Next() {
        var rev Review
        rows.Scan(
            &rev.ID, &rev.StartupID,
            &rev.ReviewerID, &rev.ReviewerName, &rev.ReviewerDesignation,
            &rev.Rating, &rev.Decision, &rev.Summary,
            &rev.Strengths, &rev.Recommendations,
            &rev.CreatedAt,
        )
        out = append(out, rev)
    }
    return out, nil
}

func (r *ReviewRepository) ListForReviewer(ctx context.Context, reviewerID string) ([]Review, error) {
    rows, err := r.pool.Query(ctx,
        `SELECT id, startup_id,
                reviewer_id, reviewer_name, reviewer_designation,
                rating, decision, summary, strengths, recommendations,
                created_at
         FROM reviews 
         WHERE reviewer_id=$1 
         ORDER BY created_at DESC`,
        reviewerID,
    )
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var out []Review
    for rows.Next() {
        var rev Review
        rows.Scan(
            &rev.ID, &rev.StartupID,
            &rev.ReviewerID, &rev.ReviewerName, &rev.ReviewerDesignation,
            &rev.Rating, &rev.Decision, &rev.Summary,
            &rev.Strengths, &rev.Recommendations,
            &rev.CreatedAt,
        )
        out = append(out, rev)
    }
    return out, nil
}