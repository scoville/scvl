package web

import (
	"net/http"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
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

func (web *Web) userRegistrationHandler(w http.ResponseWriter, r *http.Request) {
	hash := mux.Vars(r)["hash"]
	if hash == "" {
		http.Error(w, "invalid request", http.StatusUnprocessableEntity)
		return
	}
	web.engine.UserRegister(&engine.RegistrationRequest{
		Hash:     hash,
		Email:    r.FormValue("hash"),
		Password: r.FormValue("password"),
	})
}

func (web *Web) loginHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := web.store.Get(r, "scvl")
	user, err := web.engine.LoginUser(&engine.LoginUserRequest{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	session.Values["user_id"] = user.ID
	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
