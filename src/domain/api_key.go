package domain

import "time"

// APIKey is a api_key
type APIKey struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
	UserID    int        `json:"user_id" gorm:"index; not null"`
	Value     string     `json:"value" gorm:"unique_index; not null"`
	Status    string     `json:"status" gorm:"default:'active'"`
}

const APIKeyStatusActive = "active"
const APIKeyStatusInvalidated = "invalidated"
