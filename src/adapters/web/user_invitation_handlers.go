package web

import (
	"net/http"

	"github.com/gorilla/context"
	"github.com/scoville/scvl/src/domain"
	"github.com/scoville/scvl/src/engine"
)

func (web *Web) userInvitationHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := context.Get(r, "user").(*domain.User)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	web.engine.InviteUser(&engine.InviteRequest{
		FromUserID: uint(user.ID),
		Email:      r.FormValue("email"),
	})
	http.Redirect(w, r, "/", http.StatusCreated)
}

func (web *Web) invitationPageHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := context.Get(r, "user").(*domain.User)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	resp := map[string]interface{}{
		"FromUserID": user.ID,
	}
	renderTemplate(w, r, "/invite.tpl", resp)
}
