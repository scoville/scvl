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

// FormatCreatedAt formats the created_at and returns it
func (e *EmailTemplate) FormatCreatedAt() string {
	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	return e.CreatedAt.In(jst).Format("2006/01/02 15:04")
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
	if len(t.BatchEmail.Emails) < number {
		return t.BatchEmail.Emails
	}
	return t.BatchEmail.Emails[:number]
}
