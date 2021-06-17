package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/scoville/scvl/src/engine"
)

// API is the entrypoint for the api package
type API struct {
	engine *engine.Engine
}

// New returns the API instance
func New(e *engine.Engine) *API {
	return &API{e}
}

// Handle handles routing
func (api *API) Handle(r *mux.Router) {
	r.Handle("/api/shorten", api.authenticate(api.shortenHandler)).Methods(http.MethodPost)
}

func (api *API) authenticate(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Query().Get("api_key")

		apiKey, err := api.engine.FindAPIKey(key)
		if err != nil {
			sendErrorJSON(w, http.StatusForbidden, "api_key is invalid")
			return
		}
		user, err := api.engine.FindUser(uint(apiKey.UserID))
		if err != nil || user.GoogleToken == "" {
			sendErrorJSON(w, http.StatusForbidden, "user not found")
			return
		}

		context.Set(r, "user", user)

		h.ServeHTTP(w, r)
	}
}

func sendJSON(w http.ResponseWriter, statusCode int, data interface{}, headers map[string]string) {
	bytes, err := json.Marshal(data)
	if err != nil {
		sendErrorJSON(w, statusCode, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	for k, v := range headers {
		w.Header().Set(k, v)
	}
	w.WriteHeader(statusCode)
	w.Write(bytes)
}

func sendErrorJSON(w http.ResponseWriter, statusCode int, mes string) {
	if 500 <= statusCode && statusCode < 600 {
		log.Println(mes)
	}
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(struct {
		Error string `json:"error"`
	}{
		Error: mes,
	})
}

var noHeader = map[string]string{}
