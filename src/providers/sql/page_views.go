package sql

import "github.com/scoville/scvl/src/domain"

const tblPageViews = "page_views"

func (c *client) CreatePageView(slug string, pv *domain.PageView) (err error) {
	var page domain.Page
	err = c.db.Where(&domain.Page{Slug: slug}).First(&page).Error
	if err != nil {
		return
	}
	pv.PageID = int(page.ID)
	return c.db.Create(pv).Error
}
