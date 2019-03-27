package web

import (
	"fmt"
	"net/http"
	"os"
	"strings"
)

func (web *Web) oauthCallbackHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "scvl")
	retrievedState, _ := session.Values["google_state"].(string)
	if retrievedState != r.URL.Query().Get("state") {
		http.Error(w, fmt.Sprintf("Invalid session state: %s", retrievedState), http.StatusUnauthorized)
		return
	}
	u, err := fetchUserInfo(r.URL.Query().Get("code"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	allowedDomain := os.Getenv("ALLOWED_DOMAIN")
	if allowedDomain != "" && !strings.HasSuffix(u.Email, "@"+allowedDomain) {
		http.Error(w, "ログインは、Scovilleアカウントである必要があります", http.StatusUnprocessableEntity)
		return
	}
	user, err := manager.findOrCreateUser(u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	session.Values["user_id"] = user.ID
	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
