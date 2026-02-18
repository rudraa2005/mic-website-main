package repository

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Content struct {
	ID               uuid.UUID       `json:"id"`
	ContentType      string          `json:"content_type"`
	Title            string          `json:"title"`
	Description      *string         `json:"description"`
	ContentData      json.RawMessage `json:"content_data"`
	ImageURL         *string         `json:"image_url"`
	OrderIndex       int             `json:"order_index"`
	IsActive         bool            `json:"is_active"`
	CreatedAt        time.Time       `json:"created_at"`
	UpdatedAt        time.Time       `json:"updated_at"`
	RegistrationLink string          `json:"registration_link"`
}

type EventContent struct {
	ID               uuid.UUID
	Title            string
	Description      *string
	EventDate        string
	ImageURL         *string
	OrderIndex       int
	RegistrationLink string
	Venue            string
	Price            string
}

type ContentRepository struct {
	db *pgxpool.Pool
}

func NewContentRepository(db *pgxpool.Pool) *ContentRepository {
	return &ContentRepository{db: db}
}

func (r *ContentRepository) GetActiveByType(
	ctx context.Context,
	contentType string,
) ([]Content, error) {

	query := `
		SELECT
			id,
			content_type,
			title,
			description,
			content_data,
			image_url,
			order_index,
			is_active,
			COALESCE(content_data->>'registration_link', '') AS registration_link,
			created_at,
			updated_at
		FROM content
		WHERE content_type = $1
		  AND is_active = true
		ORDER BY order_index ASC
	`

	rows, err := r.db.Query(ctx, query, contentType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []Content

	for rows.Next() {
		var c Content
		err := rows.Scan(
			&c.ID,
			&c.ContentType,
			&c.Title,
			&c.Description,
			&c.ContentData,
			&c.ImageURL,
			&c.OrderIndex,
			&c.IsActive,
			&c.RegistrationLink,
			&c.CreatedAt,
			&c.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, c)
	}

	return items, nil
}

func (r *ContentRepository) GetTopResources(
	ctx context.Context,
	limit int,
) ([]Content, error) {

	query := `
		SELECT
			id,
			content_type,
			title,
			description,
			content_data,
			image_url,
			order_index,
			is_active,
			COALESCE(content_data->>'registration_link', '') AS registration_link,
			created_at,
			updated_at
		FROM content
		WHERE content_type = 'resource'
		  AND is_active = true
		ORDER BY order_index ASC
		LIMIT $1
	`

	rows, err := r.db.Query(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []Content

	for rows.Next() {
		var c Content
		err := rows.Scan(
			&c.ID,
			&c.ContentType,
			&c.Title,
			&c.Description,
			&c.ContentData,
			&c.ImageURL,
			&c.OrderIndex,
			&c.IsActive,
			&c.RegistrationLink,
			&c.CreatedAt,
			&c.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, c)
	}

	return items, nil
}

func (r *ContentRepository) GetUpcomingEvents(
	ctx context.Context,
	limit int,
) ([]EventContent, error) {

	query := `
		SELECT
  			id,
  			title,
  			description,
  			COALESCE(content_data->>'event_date', '') AS event_date,
  			image_url,
  			order_index,
  			COALESCE(content_data->>'registration_link', '') AS registration_link,
  			COALESCE(content_data->>'venue', '') AS venue,
  			COALESCE(content_data->>'price', 'Free') AS price
		FROM content
		WHERE content_type = 'event'
  			AND is_active = true
		ORDER BY 
			CASE 
				WHEN content_data->>'event_date' ~ '^\d{4}-\d{2}-\d{2}$' 
				THEN (content_data->>'event_date')::date 
				ELSE '9999-12-31'::date 
			END ASC
		LIMIT $1;
	`

	rows, err := r.db.Query(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []EventContent

	for rows.Next() {
		var e EventContent
		err := rows.Scan(
			&e.ID,
			&e.Title,
			&e.Description,
			&e.EventDate,
			&e.ImageURL,
			&e.OrderIndex,
			&e.RegistrationLink,
			&e.Venue,
			&e.Price,
		)
		if err != nil {
			return nil, err
		}
		events = append(events, e)
	}

	return events, nil
}

func (r *ContentRepository) GetAllEvents(ctx context.Context) ([]EventContent, error) {
	query := `
		SELECT
  			id,
  			title,
  			description,
  			COALESCE(content_data->>'event_date', '') AS event_date,
  			image_url,
  			order_index,
  			COALESCE(content_data->>'registration_link', '') AS registration_link,
  			COALESCE(content_data->>'venue', '') AS venue,
  			COALESCE(content_data->>'price', 'Free') AS price
		FROM content
		WHERE content_type = 'event'
  			AND is_active = true
		ORDER BY 
			CASE 
				WHEN content_data->>'event_date' ~ '^\d{4}-\d{2}-\d{2}$' 
				THEN (content_data->>'event_date')::date 
				ELSE '9999-12-31'::date 
			END ASC;
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []EventContent

	for rows.Next() {
		var e EventContent
		err := rows.Scan(
			&e.ID,
			&e.Title,
			&e.Description,
			&e.EventDate,
			&e.ImageURL,
			&e.OrderIndex,
			&e.RegistrationLink,
			&e.Venue,
			&e.Price,
		)
		if err != nil {
			return nil, err
		}
		events = append(events, e)
	}

	return events, nil
}

func (r *ContentRepository) GetAll(ctx context.Context) ([]Content, error) {
	rows, err := r.db.Query(ctx, `
		SELECT
			id,
			content_type,
			title,
			description,
			content_data,
			image_url,
			order_index,
			is_active,
			COALESCE(content_data->>'registration_link',''),
			created_at,
			updated_at
		FROM content
		ORDER BY content_type, order_index ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Content
	for rows.Next() {
		var c Content
		if err := rows.Scan(
			&c.ID,
			&c.ContentType,
			&c.Title,
			&c.Description,
			&c.ContentData,
			&c.ImageURL,
			&c.OrderIndex,
			&c.IsActive,
			&c.RegistrationLink,
			&c.CreatedAt,
			&c.UpdatedAt,
		); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, nil
}

func (r *ContentRepository) Create(ctx context.Context, c *Content) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO content (content_type, title, description, content_data, image_url, order_index, is_active)
		VALUES ($1,$2,$3,$4,$5,$6,$7)
	`,
		c.ContentType,
		c.Title,
		c.Description,
		c.ContentData,
		c.ImageURL,
		c.OrderIndex,
		c.IsActive,
	)
	log.Println(err)
	return err
}

func (r *ContentRepository) Update(ctx context.Context, id uuid.UUID, c *Content) error {
	_, err := r.db.Exec(ctx, `
		UPDATE content
		SET title=$1, description=$2, content_data=$3, image_url=$4, order_index=$5, is_active=$6
		WHERE id=$7
	`,
		c.Title,
		c.Description,
		c.ContentData,
		c.ImageURL,
		c.OrderIndex,
		c.IsActive,
		id,
	)
	return err
}

func (r *ContentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, `DELETE FROM content WHERE id=$1`, id)
	return err
}
