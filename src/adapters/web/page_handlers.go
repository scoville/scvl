package web

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-playground/form"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/scoville/scvl/src/domain"
	"github.com/scoville/scvl/src/engine"
	qrcode "github.com/skip2/go-qrcode"
	"github.com/tomasen/realip"
)

func (web *Web) pagesHandler(w http.ResponseWriter, r *http.Request) {
	bytes, _ := getFlash(w, r, "url_slug")
	resp := map[string]interface{}{}
	if bytes != nil {
		json.Unmarshal(bytes, &resp)
	}
	user, ok := context.Get(r, "user").(*domain.User)
	if !ok {
		user = &domain.User{}
		// http.Error(w, "Unauthorized", http.StatusUnauthorized)
		// return
	}
	resp["User"] = user
	var req engine.FindPagesRequest
	err := form.NewDecoder().Decode(&req, r.URL.Query())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	req.UserID = user.ID
	pages, count, err := web.engine.FindPages(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	resp["Pages"] = pages
	resp["Count"] = count
	start := req.Offset + 1
	if start < 1 {
		start = 1
	}
	resp["Start"] = start
	end := req.Offset + req.Limit
	if end > uint(count) {
		end = uint(count)
	}
	resp["End"] = end
	if req.Offset > 0 {
		prevOffset := req.Offset - req.Limit
		if prevOffset < 0 {
			prevOffset = 0
		}
		resp["PrevURL"] = fmt.Sprintf("/pages?offset=%d&limit=%d", prevOffset, req.Limit)
	}
	if req.Offset+req.Limit < uint(count) {
		nextOffset := req.Offset + req.Limit
		resp["NextURL"] = fmt.Sprintf("/pages?offset=%d&limit=%d", nextOffset, req.Limit)
	}
	loginURL, ok := context.Get(r, "login_url").(string)
	if ok {
		resp["LoginURL"] = loginURL
	}
	renderTemplate(w, r, "/pages.tpl", resp)
}

func (web *Web) shortenHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := context.Get(r, "user").(*domain.User)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	page, err := web.engine.Shorten(&engine.ShortenRequest{
		UserID:       int(user.ID),
		URL:          r.FormValue("url"),
		CustomizeOGP: r.FormValue("ogp") == "on",
		Description:  r.FormValue("description"),
		Image:        r.FormValue("image"),
		Title:        r.FormValue("title"),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	bytes, _ := json.Marshal(map[string]string{
		"URL":  page.URL,
		"Slug": page.Slug,
	})
	setFlash(w, "url_slug", bytes)
	http.Redirect(w, r, "/pages", http.StatusSeeOther)
}

func (web *Web) redirectHandler(w http.ResponseWriter, r *http.Request) {
	url, ogp, err := web.engine.Access(&engine.AccessRequest{
		Slug:      mux.Vars(r)["slug"],
		RealIP:    realip.RealIP(r),
		Referer:   r.Referer(),
		UserAgent: r.UserAgent(),
	})
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "The URL you are looking for is not found.", http.StatusNotFound)
		return
	}
	if ogp != nil {
		data := map[string]interface{}{
			"URL": url,
			"OGP": ogp,
		}
		tpl := findTemplateWithoutBase("/redirect.tpl")
		tpl.Execute(w, data)
		return
	}
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (web *Web) qrHandler(w http.ResponseWriter, r *http.Request) {
	png, err := qrcode.Encode(strings.Split(r.RequestURI, "/qr.png")[0], qrcode.Medium, 256)
	if err != nil {
		log.Println("Failed to generate QR code: ", err)
		http.Error(w, "Failed to generate QR code", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Length", strconv.Itoa(len(png)))
	if _, err := w.Write(png); err != nil {
		log.Println("Unable to write image: ", err)
		http.Error(w, "Unable to write image", http.StatusInternalServerError)
	}
}

func (web *Web) editHandler(w http.ResponseWriter, r *http.Request) {
	bytes, _ := getFlash(w, r, "message")
	resp := map[string]interface{}{}
	if bytes != nil {
		json.Unmarshal(bytes, &resp)
	}

	user, ok := context.Get(r, "user").(*domain.User)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	slug := mux.Vars(r)["slug"]
	page, err := web.engine.FindPage(slug, int(user.ID))
	if err != nil {
		http.Error(w, "The page you are looking for is not found.", http.StatusNotFound)
		return
	}

	resp["Page"] = page
	if page.OGP != nil {
		resp["OGP"] = true
	}
	renderTemplate(w, r, "/edit.tpl", resp)
}

func (web *Web) updateHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := context.Get(r, "user").(*domain.User)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	page, err := web.engine.UpdatePage(&engine.UpdatePageRequest{
		Slug:         mux.Vars(r)["slug"],
		UserID:       int(user.ID),
		URL:          r.FormValue("url"),
		CustomizeOGP: r.FormValue("ogp") == "on",
		Description:  r.FormValue("description"),
		Image:        r.FormValue("image"),
		Title:        r.FormValue("title"),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	bytes, _ := json.Marshal(map[string]string{
		"Success": "Update succeeded.",
	})
	setFlash(w, "message", bytes)
	http.Redirect(w, r, "/"+page.Slug+"/edit", http.StatusSeeOther)
}
