package sql

func (c *client) findPageBySlug(slug string) (page Page, err error) {
	err = m.db.Table("pages").
		Preload("OGP").
		First(&page, "slug = ?", slug).Error
	return
}

func (c *client) createPage(userID uint, slug, url string) (page Page, err error) {
	page = Page{
		UserID: int(userID),
		Slug:   slug,
		URL:    url,
	}
	err = m.db.Create(&page).Error
	return
}

func (c *client) updatePage(id uint, url string) (err error) {
	return m.db.Table("pages").
		Where("id = ?", id).
		Update(&Page{
			URL: url,
		}).Error
}
