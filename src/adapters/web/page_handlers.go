package web

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/scoville/scvl/src/domain"
)

func (web *Web) rootHandler(w http.ResponseWriter, r *http.Request) {
	bytes, _ := getFlash(w, r, "url_slug")
	resp := map[string]interface{}{}
	if bytes != nil {
		json.Unmarshal(bytes, &resp)
	}
	user, ok := context.Get(r, "user").(*domain.User)
	if ok {
		manager.setPagesToUser(user)
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
	url := r.FormValue("url")
	if url == "" {
		http.Error(w, "url cannot be empty", http.StatusUnprocessableEntity)
		return
	}

	slug := generateSlug(4)
	page, err := manager.createPage(user.ID, slug, url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.FormValue("ogp") == "on" {
		ogp := OGP{
			PageID:      int(page.ID),
			Description: r.FormValue("description"),
			Image:       r.FormValue("image"),
			Title:       r.FormValue("title"),
		}
		err = manager.createOGP(&ogp)
		client.SetOGPID(page.Slug, int(ogp.ID))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	client.SetURL(slug, url)
	bytes, _ := json.Marshal(map[string]string{
		"URL":  url,
		"Slug": slug,
	})
	setFlash(w, "url_slug", bytes)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (web *Web) redirectHandler(w http.ResponseWriter, r *http.Request) {
	slug := mux.Vars(r)["slug"]
	url := client.GetURL(slug)
	var ogp *OGP
	if url == "" {
		// Redisでページが見つからなかった場合の処理
		page, err := manager.findPageBySlug(slug)
		if err != nil {
			http.Error(w, "The URL you are looking for is not found.", http.StatusNotFound)
			return
		}
		url = page.URL
		client.SetURL(slug, url)
		if page.OGP != nil {
			ogp = page.OGP
			client.SetOGPID(slug, int(page.OGP.ID))
		}
	}
	ua := user_agent.New(r.UserAgent())
	if !ua.Bot() {
		name, _ := ua.Browser()
		manager.createPageView(slug, PageView{
			RealIP:      realip.RealIP(r),
			Referer:     r.Referer(),
			Mobile:      ua.Mobile(),
			Platform:    ua.Platform(),
			OS:          ua.OS(),
			BrowserName: name,
		})
	}
	var ogpID int
	if ogp == nil {
		ogpID = client.GetOGPID(slug)
	}
	if ogpID != 0 || ogp != nil {
		if ogp == nil {
			ogp, _ = manager.findOGPByID(ogpID)
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
	page, err := manager.findPageBySlug(slug)
	if err != nil {
		http.Error(w, "The page you are looking for is not found.", http.StatusNotFound)
		return
	}

	if page.UserID != int(user.ID) {
		http.Error(w, "You don't have permission to edit it.", http.StatusUnauthorized)
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
