package sql

import "github.com/scoville/scvl/src/domain"

const tblPages = "pages"

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
