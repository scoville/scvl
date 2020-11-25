package web

import (
	"log"
	"net/http"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/scoville/scvl/src/domain"
	"github.com/scoville/scvl/src/engine"
)

func (web *Web) emailsHandler(w http.ResponseWriter, r *http.Request) {
	resp := map[string]interface{}{}
	user, ok := context.Get(r, "user").(*domain.User)
	if ok {
		resp["User"] = user
		resp["SenderName"] = user.Name
		resp["BCCAddress"] = user.Email
	}
	loginURL, ok := context.Get(r, "login_url").(string)
	if ok {
		resp["LoginURL"] = loginURL
	}
	renderTemplate(w, r, "/emails.tpl", resp)
}

func (web *Web) emailCreateHandler(w http.ResponseWriter, r *http.Request) {
	var req engine.CreateEmailRequest
	req.IsPreview = r.FormValue("preview") == "1"
	req.SpreadsheetURL = r.FormValue("spreadsheet_url")
	req.SheetName = r.FormValue("sheet_name")
	req.Sender = r.FormValue("sender")
	req.Title = r.FormValue("title")
	req.Template = r.FormValue("template")
	user, ok := context.Get(r, "user").(*domain.User)
	if !ok {
		log.Println("failed to get user from the context")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	req.User = user

	emailTemplate, err := web.engine.CreateEmail(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := map[string]interface{}{
		"EmailTemplate": emailTemplate,
	}
	renderTemplate(w, r, "/email_preview.tpl", resp)
}

func (web *Web) emailSendHandler(w http.ResponseWriter, r *http.Request) {
	var req engine.CreateEmailRequest
	req.SpreadsheetURL = r.FormValue("spreadsheet_url")
	req.SheetName = r.FormValue("sheet_name")
	req.Sender = r.FormValue("sender")
	req.Title = r.FormValue("title")
	req.Template = r.FormValue("template")
	user, ok := context.Get(r, "user").(*domain.User)
	if !ok {
		log.Println("failed to get user from the context")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	req.User = user
	err := web.engine.SendEmail(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (web *Web) emailReadHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	err := web.engine.ReadEmail(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-type", "image/gif")
	w.WriteHeader(http.StatusOK)
	// 透過gifを返す
	w.Write([]byte("data:image/gif;base64,R0lGODlhAQABAIAAAAAAAP///yH5BAEAAAAALAAAAAABAAEAAAIBRAA7"))
}
