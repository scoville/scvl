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
	Title      string      `json:"title"`
	Body       string      `json:"body" gorm:"not null"`
	BatchEmail *BatchEmail `json:"batch_email"`
}

// SendingNumber Caluculates number of emails
func (t *EmailTemplate) SendingNumber() int {
	if t.BatchEmail == nil {
		return 0
	}
	return len(t.BatchEmail.Emails)
}

// FilterEmailNum return filtered emails
func (t *EmailTemplate) FilterEmailNum(number int) []*Email {
	return t.BatchEmail.Emails[:number]
}
