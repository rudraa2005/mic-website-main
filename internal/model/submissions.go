package model

import (
	"time"
)

type Submission struct {
	SubmissionID string    `json:"submission_id"`
	UserID       string    `json:"user_id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	FilePath     *string   `json:"file_path"`
	Status       string    `json:"status"`
	Stage        string    `json:"stage"`
	CompanyID    *string   `json:"company_id"`
	CompanyName  *string   `json:"company_name"`
	CompanyLogo  *string   `json:"company_logo"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
