package domain

import (
	"time"
)

// BatchEmail is the struct
type BatchEmail struct {
	ID              uint `gorm:"primary_key"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	EmailTemplateID int      `json:"email_template_id" gorm:"index; not null"`
	Sender          string   `json:"sender"`
	SpreadsheetURL  string   `json:"spreadsheet_url"`
	SheetName       string   `json:"sheet_name"`
	Emails          []*Email `json:"emails"`
	OpenCount       int      `json:"open_count"`
	SentCount       int      `json:"sent_count"`
}
