package main

import "github.com/jinzhu/gorm"

// User is a user
type User struct {
	gorm.Model
	Name  string  `json:"name"`
	Email string  `json:"email" gorm:"type:varchar(100);unique_index"`
	Pages []*Page `json:"pages"`
}

// Page is a page
type Page struct {
	gorm.Model
	UserID    int        `json:"user_id" gorm:"index; not null"`
	Slug      string     `json:"slug" gorm:"unique_index; not null"`
	URL       string     `json:"url"`
	Views     []PageView `json:"views"`
	ViewCount uint       `json:"view_count" gorm:"-"`
}

// OGP is the struct
type OGP struct {
	CanonicalURL string  `json:"canonical_url"`
	Description  string  `json:"description"`
	Image        *string `json:"image"`
	Title        string  `json:"title"`
	URL          string  `json:"url"`
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
