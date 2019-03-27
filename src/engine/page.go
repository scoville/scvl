package engine

import (
	"errors"

	"github.com/mssola/user_agent"
	"github.com/scoville/scvl/src/domain"
)

// FindPage returns the page
func (e *Engine) FindPage(slug string, userID int) (page *domain.Page, err error) {
	page, err = e.sqlClient.FindPageBySlug(slug)
	if err != nil {
		return
	}
	if page.UserID != userID {
		err = errors.New("You don't have permission to edit it.")
		return
	}
	return
}

// ShortenRequest is the request struct for Shorten function
type ShortenRequest struct {
	UserID       int
	URL          string
	CustomizeOGP bool
	Description  string
	Image        string
	Title        string
}

// Shorten shorten url
func (e *Engine) Shorten(req *ShortenRequest) (page *domain.Page, err error) {
	if req.URL == "" {
		err = errors.New("url cannot be empty")
		return
	}
	page = &domain.Page{
		UserID: req.UserID,
		Slug:   domain.GenerateSlug(5),
		URL:    req.URL,
	}
	err = e.sqlClient.CreatePage(page)
	if err != nil {
		return
	}

	e.redisClient.SetURL(page.Slug, page.URL)

	if !req.CustomizeOGP {
		return
	}
	page.OGP = &domain.OGP{
		PageID:      int(page.ID),
		Description: req.Description,
		Image:       req.Image,
		Title:       req.Title,
	}
	err = e.sqlClient.CreateOGP(page.OGP)
	if err != nil {
		return
	}

	e.redisClient.SetOGPID(page.Slug, int(page.OGP.ID))
	return
}

// AccessRequest is the request struct for Access function
type AccessRequest struct {
	Slug      string
	RealIP    string
	Referer   string
	UserAgent string
}

// Access accesses to the shorten url
func (e *Engine) Access(req *AccessRequest) (url string, ogp *domain.OGP, err error) {
	url = e.redisClient.GetURL(req.Slug)
	if url == "" {
		// Redisでページが見つからなかった場合の処理
		var page *domain.Page
		page, err = e.sqlClient.FindPageBySlug(req.Slug)
		if err != nil {
			return
		}
		url = page.URL
		e.redisClient.SetURL(page.Slug, url)
		if page.OGP != nil {
			ogp = page.OGP
			e.redisClient.SetOGPID(page.Slug, int(page.OGP.ID))
		}
	}
	ua := user_agent.New(req.UserAgent)
	if !ua.Bot() {
		name, _ := ua.Browser()
		e.sqlClient.CreatePageView(req.Slug, &domain.PageView{
			RealIP:      req.RealIP,
			Referer:     req.Referer,
			Mobile:      ua.Mobile(),
			Platform:    ua.Platform(),
			OS:          ua.OS(),
			BrowserName: name,
		})
	}
	if ogp != nil {
		return
	}
	ogpID := e.redisClient.GetOGPID(req.Slug)
	if ogpID != 0 {
		ogp, _ = e.sqlClient.FindOGPByID(ogpID)
	}
	return
}
