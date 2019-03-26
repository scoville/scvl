package domain

import (
	"time"

	"github.com/jinzhu/gorm"
)

// User is a user
type User struct {
	gorm.Model
	Name  string  `json:"name"`
	Email string  `json:"email" gorm:"type:varchar(100);unique_index"`
	Pages []*Page `json:"pages"`
	Files []*File `json:"files"`
}

// Page is a page
type Page struct {
	gorm.Model
	UserID    int        `json:"user_id" gorm:"index; not null"`
	Slug      string     `json:"slug" gorm:"unique_index; not null"`
	URL       string     `json:"url"`
	Views     []PageView `json:"views"`
	ViewCount uint       `json:"view_count" gorm:"-"`

	OGP *OGP `json:"ogp"`
}

// OGP is the struct
type OGP struct {
	gorm.Model
	PageID      int    `json:"page_id" gorm:"index; not null"`
	Description string `json:"description"`
	Image       string `json:"image"`
	Title       string `json:"title"`
}

// PageView is a pageview
type PageView struct {
	gorm.Model
	PageID      int    `json:"page_id" gorm:"index"`
	RealIP      string `json:"real_ip" gorm:"index"`
	Referer     string `json:"referer" gorm:"index"`
	Mobile      bool   `json:"mobile"`
	Platform    string `json:"platform"`
	OS          string `json:"os"`
	BrowserName string `json:"browser_name"`
}

// File is a file
type File struct {
	gorm.Model
	UserID            int        `json:"user_id" gorm:"index; not null"`
	EncryptedPassword string     `json:"encrypted_password"`
	Slug              string     `json:"slug" gorm:"unique_index; not null"`
	Deadline          *time.Time `json:"deadline"`

	Downloads     []FileDownload `json:"file_downloads"`
	DownloadCount uint           `json:"download_count" gorm:"-"`
}

// FileDownload is a file download
type FileDownload struct {
	gorm.Model
	FileID      int    `json:"file_id" gorm:"index; not null"`
	RealIP      string `json:"real_ip" gorm:"index"`
	Referer     string `json:"referer" gorm:"index"`
	Mobile      bool   `json:"mobile"`
	Platform    string `json:"platform"`
	OS          string `json:"os"`
	BrowserName string `json:"browser_name"`
}
