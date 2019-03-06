package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/mssola/user_agent"
	qrcode "github.com/skip2/go-qrcode"
	"github.com/tomasen/realip"
)

var client *redisClient
var store *sessions.CookieStore

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))

	client, err = newRedisClient()
	if err != nil {
		log.Fatalf("Failed to create redisClient: %v", err)
	}
	defer client.Close()
	setupRoutes()
	setupGoogleConfig()
	setupManager()
	defer manager.db.Close()
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func setupRoutes() {
	r := mux.NewRouter()

	r.Handle("/shorten", authenticate(shortenHandler)).Methods("POST")
	r.Handle("/", authenticate(rootHandler)).Methods("GET")

	r.HandleFunc("/{slug}/qr.png", qrHandler).Methods("GET")
	r.HandleFunc("/{slug}", redirectHandler).Methods("GET")
	r.HandleFunc("/oauth/google/callback", oauthCallbackHandler).Methods("GET")
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css/"))))
	http.Handle("/", r)
}

func authenticate(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "scvl")
		userID, ok := session.Values["user_id"].(uint)
		if ok {
			user, err := manager.findUser(userID)
			if err != nil {
				ok = false
			} else {
				context.Set(r, "user", &user)
			}
		}
		if !ok {
			state := generateSlug() + generateSlug()
			session.Values["google_state"] = state
			session.Save(r, w)
			context.Set(r, "login_url", googleConfig.AuthCodeURL(state))
		}
		h.ServeHTTP(w, r)
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	bytes, _ := getFlash(w, r, "url_slug")
	resp := map[string]interface{}{}
	if bytes != nil {
		json.Unmarshal(bytes, resp)
	}
	user, ok := context.Get(r, "user").(*User)
	if ok {
		manager.setPagesToUser(user)
		resp["User"] = user
	}
	loginURL, ok := context.Get(r, "login_url").(string)
	if ok {
		resp["LoginURL"] = loginURL
	}
	renderTemplate(w, r, "/index.tpl", resp)
}

func oauthCallbackHandler(w http.ResponseWriter, r *http.Request) {
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

func shortenHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := context.Get(r, "user").(*User)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	url := r.FormValue("url")
	slug := generateSlug()
	err := manager.createPage(user.ID, slug, url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	client.SetURL(slug, url)
	bytes, _ := json.Marshal(map[string]string{
		"URL":  url,
		"Slug": slug,
	})
	setFlash(w, "url_slug", bytes)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	slug := mux.Vars(r)["slug"]
	url := client.GetURL(slug)
	if url == "" {
		// Redisでページが見つからなかった場合の処理
		page, err := manager.findPageBySlug(slug)
		if err != nil {
			http.Error(w, "The URL you are looking for is not found.", http.StatusNotFound)
			return
		}
		url = page.URL
		client.SetURL(slug, url)
	}
	ua := user_agent.New(r.UserAgent())
	if !ua.Bot() {
		name, _ := ua.Browser()
		manager.createPageView(slug, PageView{
			RealIP:      realip.RealIP(r),
			Referer:     r.Referer(),
			Mobile:      ua.Mobile(),
			Platform:    ua.Platform(),
			OS:          ua.OS(),
			BrowserName: name,
		})
	}
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func qrHandler(w http.ResponseWriter, r *http.Request) {
	png, err := qrcode.Encode(strings.Split(r.RequestURI, "/qr.png")[0], qrcode.Medium, 256)
	if err != nil {
		log.Println("Failed to generate QR code: ", err)
		http.Error(w, "Failed to generate QR code", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Length", strconv.Itoa(len(png)))
	if _, err := w.Write(png); err != nil {
		log.Println("Unable to write image: ", err)
		http.Error(w, "Unable to write image", http.StatusInternalServerError)
	}
}

func setFlash(w http.ResponseWriter, name string, value []byte) {
	c := &http.Cookie{Name: name, Value: base64.URLEncoding.EncodeToString(value)}
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
	dc := &http.Cookie{Name: name, MaxAge: -1, Expires: time.Unix(1, 0)}
	http.SetCookie(w, dc)
	return value, nil
}
