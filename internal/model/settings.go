package model

import "time"

type Settings struct {
	UserID string `json:"user_id"`

	Theme string `json:"theme"`

	EmailNotifications bool `json:"email_notifications"`
	FeedbackAlerts     bool `json:"feedback_alerts"`
	ApplicationUpdates bool `json:"application_updates"`
	Newsletter         bool `json:"newsletter"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
