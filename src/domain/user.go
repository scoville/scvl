package domain

import (
	"fmt"
	"time"
)

// User is a user
type User struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
	Name      string     `json:"name"`
	Email     string     `json:"email" gorm:"type:varchar(100);unique_index"`

	Files             []*File  `json:"files"`
	Images            []*Image `json:"images"`
	GoogleToken       string   `json:"google_token"`
	EncryptedPassword string   `json:"-"`
	Status            string   `json:"status"`
}

// Status for User
const (
	UserStatusTemp    = "temp"
	UserStatusValid   = "valid"
	UserStatusDeleted = "deleted"
)

// SetPassword sets the password
func (w *User) SetPassword(pass string) error {
	if len(pass) < 6 {
		return fmt.Errorf("password should be greater or equal than 6 characters")
	}
	w.EncryptedPassword = Encrypt(pass)
	return nil
}

// BeforeSave is called before it is saved to the database
func (w *User) BeforeSave() error {
	if w.GoogleToken == "" && w.EncryptedPassword == "" {
		return fmt.Errorf("password is required")
	}
	return nil
}
