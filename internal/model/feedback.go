package model

import "time"

type Feedback struct {
	FeedbackID   string `json:"feedback_id"`
	SubmissionID string `json:"submission_id"`

	FacultyID    string `json:"faculty_id"`
	FacultyName  string `json:"faculty_name"`
	FacultyTitle string `json:"faculty_title"`
	FacultyField string `json:"faculty_field"`

	OverallFeedback string   `json:"overall_feedback"`
	Strengths       []string `json:"strengths"`
	Recommendations []string `json:"recommendations"`

	Rating float32 `json:"rating"`
	Status string  `json:"status"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
