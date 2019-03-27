package web

import (
	"fmt"
	"net/http"
)

func (web *Web) oauthCallbackHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := web.store.Get(r, "scvl")
	retrievedState, _ := session.Values["google_state"].(string)
	if retrievedState != r.URL.Query().Get("state") {
		http.Error(w, fmt.Sprintf("Invalid session state: %s", retrievedState), http.StatusUnauthorized)
		return
	}
	code := r.URL.Query().Get("code")
	user, err := web.engine.FindOrCreateUserByGoogleCode(code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	session.Values["user_id"] = user.ID
	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
