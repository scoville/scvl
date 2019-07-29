package sql

import (
	"github.com/scoville/scvl/src/domain"
	"github.com/scoville/scvl/src/engine"
)

const tblPages = "pages"

func (c *client) FindPages(params *engine.FindPagesRequest) (pages []*domain.Page, count int, err error) {
	db := c.db.Table(tblPages).
		Where("user_id = ?", params.UserID).
		Count(&count)

	if params.Limit != 0 {
		db = db.Limit(params.Limit)
	}
	if params.Offset != 0 {
		db = db.Offset(params.Offset)
	}

	err = db.Order("created_at DESC").
		Preload("OGP").
		Find(&pages).Error
	if err != nil {
		return
	}
	for _, p := range pages {
		c.db.Table(tblPageViews).Where("page_id = ?", p.ID).Count(&(p.ViewCount))
	}

	return
}

func (c *client) FindPageBySlug(slug string) (page *domain.Page, err error) {
	page = &domain.Page{}
	err = c.db.Table(tblPages).
		Preload("OGP").
		First(page, "slug = ?", slug).Error
	return
}

func (c *client) CreatePage(page *domain.Page) (err error) {
	err = c.db.Create(page).Error
	return
}

func (c *client) UpdatePage(page, params *domain.Page) error {
	return c.db.Table(tblPages).
		Model(page).
		Update(params).Error
}
