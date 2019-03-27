package domain

import "time"

// FileDownload is a file download
type FileDownload struct {
	ID          uint `gorm:"primary_key"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time `sql:"index"`
	FileID      int        `json:"file_id" gorm:"index; not null"`
	RealIP      string     `json:"real_ip" gorm:"index"`
	Referer     string     `json:"referer" gorm:"index"`
	Mobile      bool       `json:"mobile"`
	Platform    string     `json:"platform"`
	OS          string     `json:"os"`
	BrowserName string     `json:"browser_name"`
}
