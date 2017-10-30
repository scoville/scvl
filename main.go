package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/mssola/user_agent"
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
	resp := map[string]interface{}{}
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
	resp := map[string]interface{}{}
	user, ok := context.Get(r, "user").(*User)
	if ok {
		resp["User"] = user
	} else {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	url := r.FormValue("url")
	slug := generateSlug()
	client.SetURL(slug, url)
	err := manager.createPage(user.ID, slug, url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	manager.setPagesToUser(user)

	loginURL, ok := context.Get(r, "login_url").(string)
	if ok {
		resp["LoginURL"] = loginURL
	}
	resp["URL"] = url
	resp["Slug"] = slug
	renderTemplate(w, r, "/index.tpl", resp)
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	slug := mux.Vars(r)["slug"]
	url := client.GetURL(slug)
	if url == "" {
		http.Error(w, "The URL you are looking for is not found.", http.StatusNotFound)
		return
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
