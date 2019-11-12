package engine

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/mssola/user_agent"
	"github.com/scoville/scvl/src/domain"
)

// FindPagesRequest is the request struct for FindPages()
type FindPagesRequest struct {
	Limit  uint `form:"limit"`
	Offset uint `form:"offset"`
	UserID uint `form:"-"`
}

// FindPages returns the pages
func (e *Engine) FindPages(req *FindPagesRequest) (pages []*domain.Page, count int, err error) {
	if req.Limit == 0 {
		req.Limit = 20
	}
	return e.sqlClient.FindPages(req)
}

// FindPage returns the page
func (e *Engine) FindPage(slug string, userID int) (page *domain.Page, err error) {
	page, err = e.sqlClient.FindPageBySlug(slug)
	if err != nil {
		return
	}
	if userID != 0 && page.UserID != userID {
		err = errors.New("You don't have permission to edit it")
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
	title, err := e.fetchTitle(req.UserID, req.URL)
	if err != nil {
		return
	}

	page = &domain.Page{
		UserID: req.UserID,
		Slug:   domain.GenerateSlug(5),
		Title:  title,
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

// UpdatePageRequest is the request struct for UpdatePage function
type UpdatePageRequest struct {
	Slug         string
	UserID       int
	URL          string
	CustomizeOGP bool
	Description  string
	Image        string
	Title        string
}

// UpdatePage updates the page
func (e *Engine) UpdatePage(req *UpdatePageRequest) (page *domain.Page, err error) {
	if req.URL == "" {
		err = errors.New("url cannot be empty")
		return
	}

	page, err = e.sqlClient.FindPageBySlug(req.Slug)
	if err != nil {
		return
	}

	if page.UserID != req.UserID {
		err = errors.New("You don't have permission to edit it")
		return
	}

	title, err := e.fetchTitle(req.UserID, req.URL)
	if err != nil {
		return
	}

	err = e.sqlClient.UpdatePage(page, &domain.Page{URL: req.URL, Title: title})
	if err != nil {
		return
	}

	e.redisClient.SetURL(req.Slug, req.URL)

	if req.CustomizeOGP {
		var ogpID int
		if page.OGP == nil {
			ogp := domain.OGP{
				PageID:      int(page.ID),
				Description: req.Description,
				Image:       req.Image,
				Title:       req.Title,
			}
			err = e.sqlClient.CreateOGP(&ogp)
			ogpID = int(ogp.ID)
		} else {
			err = e.sqlClient.UpdateOGP(page.OGP.ID, &domain.OGP{
				Description: req.Description,
				Image:       req.Image,
				Title:       req.Title,
			})
			ogpID = int(page.OGP.ID)
		}
		if err != nil {
			return
		}
		e.redisClient.SetOGPID(page.Slug, ogpID)
	} else if page.OGP != nil {
		e.redisClient.DeleteOGPID(page.Slug)
		err = e.sqlClient.DeleteOGP(page.OGP.ID)
		if err != nil {
			return
		}
	}
	return
}

func (e *Engine) fetchTitle(userID int, url string) (title string, err error) {
	user, err := e.sqlClient.FindUser(uint(userID))
	if err != nil {
		return
	}
	if strings.HasPrefix(url, "https://docs.google.com/") && !strings.HasPrefix(url, "https://docs.google.com/forms/") {
		paths := strings.Split(url, "/")
		if len(paths) > 5 {
			title, err = e.googleClient.GetDriveFileTitle(user, paths[5])
		}
		if err == nil {
			return
		}
	}
	var resp *http.Response
	resp, err = http.Get("https://ogp.en-courage.com?url=" + url)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	var ogp domain.OGP
	err = json.NewDecoder(resp.Body).Decode(&ogp)
	if err != nil {
		return
	}
	title = ogp.Title
	return
}
