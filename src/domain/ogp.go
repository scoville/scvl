package domain

import "time"

// OGP is the struct
type OGP struct {
	ID          uint `gorm:"primary_key"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time `sql:"index"`
	PageID      int        `json:"page_id" gorm:"index; not null"`
	Description string     `json:"description"`
	Image       string     `json:"image"`
	Title       string     `json:"title"`
}
