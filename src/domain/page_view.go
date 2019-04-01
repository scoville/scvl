package domain

import "time"

// PageView is a pageview
type PageView struct {
	ID          uint `gorm:"primary_key"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time `sql:"index"`
	PageID      int        `json:"page_id" gorm:"index"`
	RealIP      string     `json:"real_ip" gorm:"index"`
	Referer     string     `json:"referer" gorm:"index"`
	Mobile      bool       `json:"mobile"`
	Platform    string     `json:"platform"`
	OS          string     `json:"os"`
	BrowserName string     `json:"browser_name"`
}
