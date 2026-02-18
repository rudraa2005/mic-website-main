package model

import "time"

type Query struct {
	QueryID    string  `json:"query_id"`
	UserID     string  `json:"user_id"`
	FacultyID  string  `json:"faculty_id"`
	FeedbackID *string `json:"feedback_id,omitempty"`
	Status     string  `json:"status"`

	QueryText string `json:"query"`
	Priority  string `json:"priority"`

	Response *string `json:"response,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
