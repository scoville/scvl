package domain

import (
	"errors"
	"path/filepath"
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
	Email         *FileEmail     `json:"email"`
}

// Name returns the name of the file
func (f *File) Name() string {
	return filepath.Base(f.Path)
}

// Downloadable returns error if it cannot be downloaded
func (f *File) Downloadable() error {
	if f.DownloadLimit > 0 && f.DownloadCount >= uint(f.DownloadLimit) {
		return errors.New("このファイルをダウンロード可能な回数が制限を超えました。ファイルのアップロード者にお問い合わせください。")
	}
	if f.Deadline != nil && f.Deadline.Sub(time.Now()) < 0 {
		return errors.New("このファイルをダウンロード可能な期限が過ぎました。ファイルのアップロード者にお問い合わせください。")
	}
	return nil
}

// FormatDeadline formats the deadline and returns it
func (f *File) FormatDeadline() string {
	if f.Deadline == nil {
		return "無期限"
	}
	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	return f.Deadline.In(jst).Format("2006/01/02 15:04")
}

// RemainingDownloadableCount returns remaining downloadable count for the file
func (f *File) RemainingDownloadableCount() int {
	return f.DownloadLimit - int(f.DownloadCount)
}
