package web

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/scoville/scvl/src/domain"
	"github.com/scoville/scvl/src/engine"
)

func (web *Web) filesHandler(w http.ResponseWriter, r *http.Request) {
	bytes, _ := getFlash(w, r, "file_slug")
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
	renderTemplate(w, r, "/files.tpl", resp)
}

func (web *Web) fileUploadHandler(w http.ResponseWriter, r *http.Request) {
	var req engine.UploadFileRequest
	var err error

	user, ok := context.Get(r, "user").(*domain.User)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	req.UserID = int(user.ID)

	fDownloadLimit := r.FormValue("download_limit")
	req.DownloadLimit, err = strconv.Atoi(fDownloadLimit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fValidDays := r.FormValue("valid_days")
	req.ValidDays, err = strconv.Atoi(fValidDays)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	req.Email = r.FormValue("email")
	req.Password = r.FormValue("password")

	f, info, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer f.Close()
	req.File = f

	req.FileName = info.Filename
	req.FileSize = info.Size

	file, err := web.engine.UploadFile(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	bytes, _ := json.Marshal(map[string]string{
		"Slug": file.Slug,
	})
	setFlash(w, "file_slug", bytes)
	http.Redirect(w, r, "/files", http.StatusSeeOther)
}

func (web *Web) fileDownloadHandler(w http.ResponseWriter, r *http.Request) {
	slug := mux.Vars(r)["slug"]
	password := r.FormValue("password")

	data, err := web.engine.DownloadFile(slug, password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(data)
}
