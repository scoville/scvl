package web

import (
	"net/http"

	"github.com/gorilla/context"
	"github.com/scoville/scvl/src/domain"
)

func (web *Web) publishAPIKeyHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := context.Get(r, "user").(*domain.User)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	err := web.engine.UpdateUserAPIKey(user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/pages", http.StatusSeeOther)
}
