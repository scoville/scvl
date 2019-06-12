package domain

import (
	"path/filepath"
	"time"
)

// Image is an image
type Image struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
	UserID    int        `json:"user_id" gorm:"index; not null"`
	URL       string     `json:"url" gorm:"unique_index; not null"`
}

// Name returns the name of the image
func (i *Image) Name() string {
	return filepath.Base(i.URL)
}
