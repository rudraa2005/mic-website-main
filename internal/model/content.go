package model

import "time"

type BaseContent struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	IsActive  bool      `json:"is_active"`
	Order     int       `json:"order"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type AboutContent struct {
	BaseContent
	Description string `json:"description"`
	Icon        string `json:"icon,omitempty"`
}

type Event struct {
	BaseContent
	Description      string    `json:"description"`
	EventDate        time.Time `json:"event_date"`
	Venue            string    `json:"venue"`
	Price            string    `json:"price"`
	ImageURL         string    `json:"image_url"`
	Status           string    `json:"status"`
	RegistrationLink string    `json:"registration_link"`
}

type Resource struct {
	BaseContent
	Description   string `json:"description"`
	Category      string `json:"category"`
	FileURL       string `json:"file_url"`
	Format        string `json:"format"`   // PDF, Video, Docs
	Duration      string `json:"duration"` // optional display
	IsFeatured    bool   `json:"is_featured"`
	DownloadCount int    `json:"download_count"`
}
