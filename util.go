package main

import (
	"math/rand"
	"net/http"
)

func generateSlug() string {
	urlChars := []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	s := make([]rune, 4)
	for i := range s {
		s[i] = urlChars[rand.Intn(len(urlChars))]
	}
	return string(s)
}

func generateGoogleState(w http.ResponseWriter, r *http.Request) string {
	session, _ := store.Get(r, "scvl")
	state := generateSlug() + generateSlug()
	session.Values["google_state"] = state
	session.Save(r, w)
	return state
}
