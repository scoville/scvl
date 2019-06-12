package web

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/context"
	"github.com/scoville/scvl/src/domain"
	"github.com/scoville/scvl/src/engine"
)

func (web *Web) imagesHandler(w http.ResponseWriter, r *http.Request) {
	bytes, _ := getFlash(w, r, "image_url")
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
	renderTemplate(w, r, "/images.tpl", resp)
}

func (web *Web) imageUploadHandler(w http.ResponseWriter, r *http.Request) {
	var req engine.UploadImageRequest
	var err error

	user, ok := context.Get(r, "user").(*domain.User)
	if !ok {
		log.Println("failed to get user from the context")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	req.UserID = int(user.ID)

	f, info, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer f.Close()
	req.File = f

	req.FileName = info.Filename

	img, err := web.engine.UploadImage(req)
	if err != nil {
		log.Println("upload image failed: " + err.Error())
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	bytes, _ := json.Marshal(map[string]string{
		"URL": img.URL,
	})
	setFlash(w, "image_url", bytes)
	http.Redirect(w, r, "/images", http.StatusSeeOther)
}
