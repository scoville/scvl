package domain

import "time"

// User is a user
type User struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
	Name      string     `json:"name"`
	Email     string     `json:"email" gorm:"type:varchar(100);unique_index"`
	APIKey    string     `json:"api_key" gorm:"unique_index"`

	Files       []*File  `json:"files"`
	Images      []*Image `json:"images"`
	GoogleToken string   `json:"google_token"`
}
