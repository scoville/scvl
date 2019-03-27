package domain

import (
	"time"
)

// File is a file
type File struct {
	ID                uint `gorm:"primary_key"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         *time.Time `sql:"index"`
	UserID            int        `json:"user_id" gorm:"index; not null"`
	EncryptedPassword string     `json:"encrypted_password"`
	Slug              string     `json:"slug" gorm:"unique_index; not null"`
	Deadline          *time.Time `json:"deadline"`
	Path              string     `json:"path"`
	DownloadLimit     int        `json:"download_limit"`

	Downloads     []FileDownload `json:"file_downloads"`
	DownloadCount uint           `json:"download_count" gorm:"-"`
}
