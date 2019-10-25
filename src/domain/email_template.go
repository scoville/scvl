package domain

import (
	"time"
)

// EmailTemplate is the struct
type EmailTemplate struct {
	ID         uint `gorm:"primary_key"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	UserID     int         `json:"user_id" gorm:"index; not null"`
	Body       string      `json:"body" gorm:"not null"`
	BatchEmail *BatchEmail `json:"batch_email"`
}
