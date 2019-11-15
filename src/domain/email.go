package domain

import (
	"errors"
	"strings"
	"time"
)

// NewEmail initializes the email and return it
func NewEmail(template, to, title string, variables map[string]string) (email *Email, err error) {
	email = &Email{}
	email.Body = template
	email.To = to
	email.Title = title
	for k, v := range variables {
		email.Body = strings.ReplaceAll(email.Body, "[["+k+"]]", v)
		email.Title = strings.ReplaceAll(email.Title, "[["+k+"]]", v)
	}
	for _, s := range []string{email.Body, email.Title} {
		if strings.Contains(s, "[[") && strings.Contains(s, "]]") {
			err = errors.New("このメールには置換されていない変数が含まれているか、または、変数以外に[[ ]]の文字が文章内で使われています。")
			return
		}
	}
	return
}

// Email is the struct
type Email struct {
	ID            uint `gorm:"primary_key"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	BatchEmailID  int        `json:"batch_email_id" gorm:"index; not null"`
	To            string     `json:"to"`
	Title         string     `json:"title"`
	Body          string     `json:"body"`
	OpenedAt      *time.Time `json:"opened_at"`
	LinkClickedAt *time.Time `json:"link_clicked_at"`
}
