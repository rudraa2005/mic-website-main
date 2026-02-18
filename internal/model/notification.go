package model

import "time"

type Notification struct {
	ID           string
	UserID       string
	Role         string
	Type         string
	Body         string
	oldStatus    string
	newStatus    string
	SubmissionID *string
	IsRead       bool
	CreatedAt    time.Time
	Title        string
}
