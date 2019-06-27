package web

import (
	"log"
	"net/http"

	"github.com/gorilla/context"
	"github.com/scoville/scvl/src/domain"
)

func (web *Web) authenticate(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := web.store.Get(r, "scvl")
		userID, ok := session.Values["user_id"].(uint)
		if ok {
			user, err := web.engine.FindUser(userID)
			if err != nil || user.GoogleToken == "" {
				log.Println("Failed to find user: " + err.Error())
				ok = false
			} else {
				context.Set(r, "user", user)
			}
		}
		if !ok {
			state := domain.GenerateSlug(8)
			session.Values["google_state"] = state
			session.Save(r, w)
			context.Set(r, "login_url", web.engine.AuthCodeURL(state))
		}
		h.ServeHTTP(w, r)
	}
}
