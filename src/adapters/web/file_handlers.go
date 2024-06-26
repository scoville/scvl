package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/go-playground/form"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/scoville/scvl/src/domain"
	"github.com/scoville/scvl/src/engine"
	"github.com/tomasen/realip"
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
		resp["SenderName"] = user.Name
		resp["BCCAddress"] = user.Email
		var req engine.FindFilesRequest
		err := form.NewDecoder().Decode(&req, r.URL.Query())
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		req.UserID = user.ID
		files, count, err := web.engine.FindFiles(&req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		resp["Files"] = files
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
		q := ""
		if req.Query != "" {
			q += "&q=" + req.Query
			resp["Query"] = req.Query
		}

		if req.Offset > 0 {
			prevOffset := int(req.Offset) - int(req.Limit)
			if prevOffset < 0 {
				prevOffset = 0
			}
			resp["PrevURL"] = fmt.Sprintf("/?offset=%d&limit=%d", prevOffset, req.Limit) + q

		}
		if req.Offset+req.Limit < uint(count) {
			nextOffset := req.Offset + req.Limit
			resp["NextURL"] = fmt.Sprintf("/?offset=%d&limit=%d", nextOffset, req.Limit) + q
		}
	} else {
		resp["Files"] = make([]*domain.File, 0)
	}
	loginURL, ok := context.Get(r, "login_url").(string)
	if ok {
		resp["LoginURL"] = loginURL
	}
	renderTemplate(w, r, "/files.html", resp)
}

func (web *Web) fileUploadHandler(w http.ResponseWriter, r *http.Request) {
	var req engine.UploadFileRequest
	var err error

	user, ok := context.Get(r, "user").(*domain.User)
	if !ok {
		log.Println("failed to get user from the context")
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

	if r.FormValue("direct_download") == "on" && req.Password != "" {
		http.Error(w, "Cannot set password when you skip the downloading page", http.StatusUnprocessableEntity)
		return
	}

	req.SendEmail = r.FormValue("email") == "on"
	req.DirectDownload = r.FormValue("direct_download") == "on"
	req.ReceiverAddress = r.FormValue("receiver_address")
	req.ReceiverName = r.FormValue("receiver_name")
	req.SenderName = r.FormValue("sender_name")
	req.BCCAddress = r.FormValue("bcc_address")
	req.Message = r.FormValue("message")
	req.SendPassword = r.FormValue("send_password") == "on"

	file, err := web.engine.UploadFile(req)
	if err != nil {
		log.Println("upload file failed: " + err.Error())
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	bytes, _ := json.Marshal(map[string]string{
		"Slug": file.Slug,
	})
	setFlash(w, "file_slug", bytes)
	http.Redirect(w, r, "/files", http.StatusSeeOther)
}

func (web *Web) fileShowHandler(w http.ResponseWriter, r *http.Request) {
	bytes, _ := getFlash(w, r, "message")
	resp := map[string]interface{}{}
	if bytes != nil {
		json.Unmarshal(bytes, &resp)
	}

	slug := mux.Vars(r)["slug"]
	file, err := web.engine.FindFile(slug, 0)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if file.DirectDownload {
		web.downloadFile(w, r, slug, "")
		return
	}

	resp["File"] = file

	if err := file.Downloadable(); err != nil {
		log.Printf("File[ID: %d] - %s", file.ID, err.Error())
		resp["Downloadable"] = false
		resp["Error"] = err.Error()
	} else {
		resp["Downloadable"] = true
	}

	renderTemplate(w, r, "/file.html", resp)
}

func (web *Web) fileDownloadHandler(w http.ResponseWriter, r *http.Request) {
	web.downloadFile(w, r, mux.Vars(r)["slug"], r.FormValue("password"))
}

func (web *Web) downloadFile(w http.ResponseWriter, r *http.Request, slug, password string) {
	fileName, data, err := web.engine.DownloadFile(&engine.DownloadFileRequest{
		Slug:      slug,
		Password:  password,
		RealIP:    realip.RealIP(r),
		Referer:   r.Referer(),
		UserAgent: r.UserAgent(),
	})
	if err != nil {
		bytes, _ := json.Marshal(map[string]string{
			"Error": "Download Failed",
		})
		log.Println("Download failed: " + err.Error())
		setFlash(w, "message", bytes)
		http.Redirect(w, r, "/files/"+slug, http.StatusSeeOther)
		return
	}
	mime := http.DetectContentType(data)
	fileSize := len(string(data))

	urlEncodedFileName := url.QueryEscape(fileName)
	contentDisposition := fmt.Sprintf("attachment; filename=\"%s\"; filename*=UTF-8''%s", fileName, urlEncodedFileName)

	w.Header().Set("Content-Type", mime)
	w.Header().Set("Content-Disposition", contentDisposition)
	w.Header().Set("Expires", "0")
	w.Header().Set("Content-Transfer-Encoding", "binary")
	w.Header().Set("Content-Length", strconv.Itoa(fileSize))
	w.Header().Set("Content-Control", "private, no-transform, no-store, must-revalidate")

	http.ServeContent(w, r, fileName, time.Now(), bytes.NewReader(data))
}

func (web *Web) editFileHandler(w http.ResponseWriter, r *http.Request) {
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
	file, err := web.engine.FindFile(slug, int(user.ID))
	if err != nil {
		http.Error(w, "The page you are looking for is not found.", http.StatusNotFound)
		return
	}

	resp["File"] = file
	renderTemplate(w, r, "/file_edit.html", resp)
}

func (web *Web) updateFileHandler(w http.ResponseWriter, r *http.Request) {
	var req engine.UpdateFileRequest
	var err error

	user, ok := context.Get(r, "user").(*domain.User)
	if !ok {
		log.Println("failed to get user from the context")
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
	req.Slug = mux.Vars(r)["slug"]

	file, err := web.engine.UpdateFile(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	bytes, _ := json.Marshal(map[string]string{
		"Success": "Update succeeded.",
	})
	setFlash(w, "message", bytes)
	http.Redirect(w, r, "/files/"+file.Slug+"/edit", http.StatusSeeOther)
}

func (web *Web) destroyFileHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := context.Get(r, "user").(*domain.User)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	_, err := web.engine.UpdateFile(&engine.UpdateFileRequest{
		Slug:   mux.Vars(r)["slug"],
		Status: domain.FileStatusDeleted,
		UserID: int(user.ID),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	bytes, _ := json.Marshal(map[string]string{
		"Success": "Successfully deleted.",
	})
	setFlash(w, "message", bytes)
	http.Redirect(w, r, "/files", http.StatusSeeOther)
}
