package model

import "time"

type Startup struct{
	ID	int64	`json:"id"`
	OwnerID string	`json:"owner_id"`
	Title	string	`json:"title"`
	Description	string	`json:"description"`
	Stage	string	`json:"stage"`
	Department	string	`json:"department"`
	CreatedAt	time.Time	`json:"created_at"`
	UpdatedAt	time.Time	`json:"updated_at"`
}

