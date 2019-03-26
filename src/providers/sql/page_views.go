package sql

const tblPageViews = "page_views"

func (c *client) createPageView(slug string, pv PageView) (err error) {
	var page Page
	err = m.db.Where(&Page{Slug: slug}).First(&page).Error
	if err != nil {
		return
	}
	pv.PageID = int(page.ID)
	return m.db.Create(&pv).Error
}
