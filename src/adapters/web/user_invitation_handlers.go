package web

import (
	"net/http"

	"github.com/gorilla/context"
	"github.com/scoville/scvl/src/domain"
	"github.com/scoville/scvl/src/engine"
)

func (web *Web) invitationCreateHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := context.Get(r, "user").(*domain.User)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	web.engine.InviteUser(&engine.InviteRequest{
		FromUserID: uint(user.ID),
		Email:      r.FormValue("email"),
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (web *Web) invitationsHandler(w http.ResponseWriter, r *http.Request) {
	_, ok := context.Get(r, "user").(*domain.User)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	resp := map[string]interface{}{}
	renderTemplate(w, r, "/invite.tpl", resp)
}
