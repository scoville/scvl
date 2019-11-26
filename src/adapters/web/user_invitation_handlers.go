package web

import (
	"net/http"
	"os"

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
	invitation, err := web.engine.InviteUser(&engine.InviteRequest{
		FromUserID: uint(user.ID),
		Email:      r.FormValue("email"),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	url := "https://" + os.Getenv("MAIN_DOMAIN") + "/register/" + invitation.Hash
	resp := map[string]interface{}{
		"InvitationURL": url,
	}
	renderTemplate(w, r, "/invitation.tpl", resp)
}

func (web *Web) invitationsHandler(w http.ResponseWriter, r *http.Request) {
	_, ok := context.Get(r, "user").(*domain.User)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	resp := map[string]interface{}{}
	renderTemplate(w, r, "/invitations.tpl", resp)
}
