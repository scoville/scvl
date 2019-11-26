package web

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/scoville/scvl/src/engine"
)

func (web *Web) userRegistrationHandler(w http.ResponseWriter, r *http.Request) {
	user, err := web.engine.UserRegister(&engine.RegistrationRequest{
		Hash:     r.FormValue("hash"),
		Password: r.FormValue("password"),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	session, _ := web.store.Get(r, "scvl")
	session.Values["user_id"] = user.ID
	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
func (web *Web) userRegistrationPageHandler(w http.ResponseWriter, r *http.Request) {
	hash := mux.Vars(r)["hash"]
	invitation, err := web.engine.FindInvitation(&engine.FindInvitationRequest{
		Hash: hash,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	resp := map[string]interface{}{
		"Email": invitation.ToUser.Email,
		"Hash":  hash,
	}
	renderTemplate(w, r, "/register.tpl", resp)
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
