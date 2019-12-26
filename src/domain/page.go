package domain

import "time"

// Page is a page
type Page struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
	UserID    int        `json:"user_id" gorm:"index; not null"`
	Slug      string     `json:"slug" gorm:"unique_index; not null"`
	Title     string     `json:"title"`
	URL       string     `json:"url"`
	Views     []PageView `json:"views"`
	ViewCount uint       `json:"view_count" gorm:"-"`

	OGP *OGP `json:"ogp"`
}

// PageSlugLength is the length for shorten url path.
const PageSlugLength = 5
