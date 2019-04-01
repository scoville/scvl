package domain

import (
	"bytes"
	"html/template"
	"strings"
	"time"
)

// Email is the struct
type Email struct {
	ID              uint `gorm:"primary_key"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       *time.Time `sql:"index"`
	FileID          int        `json:"file_id" gorm:"index; not null"`
	SenderName      string     `json:"sender_name"`
	ReceiverName    string     `json:"receiver_name"`
	ReceiverAddress string     `json:"receiver_address"`
	BCCAddress      string     `json:"bcc_address"`
	Message         string     `json:"message"`
}

// FormatCreatedAt formats the created_at and returns it
func (e *Email) FormatCreatedAt() string {
	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	return e.CreatedAt.In(jst).Format("2006/01/02 15:04")
}

// FormatMessage formats the message and returns it
func (e *Email) FormatMessage() template.HTML {
	return template.HTML(strings.Replace(template.HTMLEscapeString(e.Message), "\n", "<br>", -1))
}

var emailTemplate *template.Template

// HTML returns email body in HTML
func (e *Email) HTML(baseURL string, file *File) (string, error) {
	if emailTemplate == nil {
		emailTemplate = template.Must(template.ParseFiles("templates/email.tpl"))
	}

	var buf bytes.Buffer
	err := emailTemplate.Execute(&buf, map[string]interface{}{
		"BaseURL": baseURL,
		"File":    file,
	})
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
