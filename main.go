package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

var client *redisClient

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	var err error
	client, err = newRedisClient()
	if err != nil {
		log.Fatalf("Failed to create redisClient: %v", err)
	}
	defer client.Close()
	setupRoutes()
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func setupRoutes() {
	r := mux.NewRouter()
	r.HandleFunc("/shorten", shortenHandler).Methods("POST")
	r.HandleFunc("/{slug}", redirectHandler).Methods("GET")
	r.HandleFunc("/", rootHandler).Methods("GET")
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css/"))))
	http.Handle("/", r)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, r, "/index.tpl", map[string]interface{}{})
}

func shortenHandler(w http.ResponseWriter, r *http.Request) {
	url := r.FormValue("url")
	slug := generateSlug()
	client.SetURL(slug, url)
	renderTemplate(w, r, "/index.tpl", map[string]interface{}{
		"URL":  url,
		"Slug": slug,
	})
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	slug := mux.Vars(r)["slug"]
	url := client.GetURL(slug)
	if url == "" {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "The URL you are looking for is not found.")
		return
	}
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}
