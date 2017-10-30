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
	UserID    int        `json:"user_id" gorm:"index"`
	Slug      string     `json:"slug" gorm:"index"`
	URL       string     `json:"url"`
	Views     []PageView `json:"views"`
	ViewCount uint       `json:"view_count" gorm:"-"`
}

// PageView is a pageview
type PageView struct {
	gorm.Model
	PageID      int    `json:"page_id" gorm:"index"`
	RemoteAddr  string `json:"remote_addr" gorm:"index"`
	Referer     string `json:"referer" gorm:"index"`
	Mobile      bool   `json:"mobile"`
	Platform    string `json:"platform"`
	OS          string `json:"os"`
	BrowserName string `json:"browser_name"`
}
