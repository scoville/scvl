package web

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/tomasen/realip"
	"github.com/skip2/go-qrcode"
	"github.com/scoville/scvl/src/domain"
	"github.com/scoville/scvl/src/engine"
)

func (web *Web) rootHandler(w http.ResponseWriter, r *http.Request) {
	bytes, _ := getFlash(w, r, "url_slug")
	resp := map[string]interface{}{}
	if bytes != nil {
		json.Unmarshal(bytes, &resp)
	}
	user, ok := context.Get(r, "user").(*domain.User)
	if ok {
		resp["User"] = user
	}
	loginURL, ok := context.Get(r, "login_url").(string)
	if ok {
		resp["LoginURL"] = loginURL
	}
	renderTemplate(w, r, "/index.tpl", resp)
}

func (web *Web) shortenHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := context.Get(r, "user").(*domain.User)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	req := &engine.ShortenRequest{
		UserID: int(user.ID),
		URL: r.FormValue("url"),
		Description: r.FormValue("description"),
		Image:       r.FormValue("image"),
		Title:       r.FormValue("title"),
	}
	if r.FormValue("ogp") == "on" {
		req.CustomizeOGP = true
	}

	page, err := web.engine.Shorten(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	bytes, _ := json.Marshal(map[string]string{
		"URL":  page.URL,
		"Slug": page.Slug,
	})
	setFlash(w, "url_slug", bytes)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (web *Web) redirectHandler(w http.ResponseWriter, r *http.Request) {
	url, ogp, err := web.engine.Access(&engine.AccessRequest{
		Slug: mux.Vars(r)["slug"],
		RealIP:      realip.RealIP(r),
		Referer:     r.Referer(),
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

	slug := mux.Vars(r)["slug"]
	page, err := manager.findPageBySlug(slug)
	if err != nil {
		http.Error(w, "The page you are looking for is not found.", http.StatusNotFound)
		return
	}

	if page.UserID != int(user.ID) {
		http.Error(w, "You don't have permission to edit it.", http.StatusUnauthorized)
		return
	}

	url := r.FormValue("url")
	if url == "" {
		http.Error(w, "url cannot be empty", http.StatusUnprocessableEntity)
		return
	}
	if err := manager.updatePage(page.ID, url); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	client.SetURL(slug, url)
	if r.FormValue("ogp") == "on" {
		var ogpID int
		if page.OGP == nil {
			ogp := OGP{
				PageID:      int(page.ID),
				Description: r.FormValue("description"),
				Image:       r.FormValue("image"),
				Title:       r.FormValue("title"),
			}
			err = manager.createOGP(&ogp)
			ogpID = int(ogp.ID)
		} else {
			err = manager.updateOGP(page.OGP.ID, OGP{
				Description: r.FormValue("description"),
				Image:       r.FormValue("image"),
				Title:       r.FormValue("title"),
			})
			ogpID = int(page.OGP.ID)
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		client.SetOGPID(page.Slug, ogpID)
	} else if page.OGP != nil {
		client.DeleteOGPID(page.Slug)
		err = manager.deleteOGP(page.OGP.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	bytes, _ := json.Marshal(map[string]string{
		"Success": "Update succeeded.",
	})
	setFlash(w, "message", bytes)
	http.Redirect(w, r, "/"+slug+"/edit", http.StatusSeeOther)
}
