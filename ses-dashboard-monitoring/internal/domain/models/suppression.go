package models

import "time"

type Suppression struct {
	ID              int       `json:"id"`
	Email           string    `json:"email"`
	Reason          string    `json:"reason"`
	Source          string    `json:"source"`
	SuppressionType string    `json:"suppression_type"`
	AWSStatus       string    `json:"aws_status"`
	IsActive        bool      `json:"is_active"`
	AddedBy         int       `json:"added_by"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
