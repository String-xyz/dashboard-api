package model

import "time"

type Platform struct {
	ID            string     `json:"id" db:"id"`
	CreatedAt     time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt     time.Time  `json:"updatedAt" db:"updated_at"`
	DeactivatedAt *time.Time `json:"deactivatedAt,omitempty" db:"deactivated_at"`
	Type          string     `json:"type" db:"type"`
	Description   string     `json:"description" db:"description"`
	ApiKey        string     `json:"apiKey" db:"api_key"`
}
