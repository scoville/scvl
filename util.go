package main

import (
	"encoding/base64"
	"math/rand"
	"net/http"
	"time"
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

func setFlash(w http.ResponseWriter, name string, value []byte) {
	c := &http.Cookie{Path: "/", Name: name, Value: base64.URLEncoding.EncodeToString(value)}
	http.SetCookie(w, c)
}

func getFlash(w http.ResponseWriter, r *http.Request, name string) ([]byte, error) {
	c, err := r.Cookie(name)
	if err != nil {
		switch err {
		case http.ErrNoCookie:
			return nil, nil
		default:
			return nil, err
		}
	}
	value, err := base64.URLEncoding.DecodeString(c.Value)
	if err != nil {
		return nil, err
	}
	dc := &http.Cookie{Path: "/", Name: name, MaxAge: -1, Expires: time.Unix(1, 0)}
	http.SetCookie(w, dc)
	return value, nil
}
