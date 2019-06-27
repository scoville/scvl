package web

import (
	"encoding/base64"
	"net/http"
	"os"
	"time"
)

func setFlash(w http.ResponseWriter, name string, value []byte) {
	c := &http.Cookie{Path: "/", Name: name, Value: base64.URLEncoding.EncodeToString(value)}
	c.Domain = os.Getenv("MAIN_DOMAIN")
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
