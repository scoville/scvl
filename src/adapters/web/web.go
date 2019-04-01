package web

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/scoville/scvl/src/engine"
)

// Digest to embed to asset
var Digest string

// Web is the entrypoint for the web package
type Web struct {
	engine *engine.Engine
	store  *sessions.CookieStore
}

// New returns the Web instance
func New(e *engine.Engine, sessionSecret string) *Web {
	return &Web{
		engine: e,
		store:  sessions.NewCookieStore([]byte(sessionSecret)),
	}
}

// Start starts listen and serve
func (web *Web) Start(port string) error {
	r := mux.NewRouter()

	r.Handle("/shorten", web.authenticate(web.shortenHandler)).Methods(http.MethodPost)
	r.Handle("/", web.authenticate(web.rootHandler)).Methods(http.MethodGet)
	r.Handle("/files", web.authenticate(web.filesHandler)).Methods(http.MethodGet)
	r.Handle("/files", web.authenticate(web.fileUploadHandler)).Methods(http.MethodPost)
	r.HandleFunc("/files/{slug}", web.fileShowHandler).Methods(http.MethodGet)
	r.Handle("/files/{slug}/download", web.authenticate(web.fileDownloadHandler)).Methods(http.MethodPost)
	r.Handle("/files/{slug}/edit", web.authenticate(web.editFileHandler)).Methods(http.MethodGet)
	r.HandleFunc("/files/{slug}", web.authenticate(web.updateFileHandler)).Methods(http.MethodPost, http.MethodPut, http.MethodPatch)

	r.HandleFunc("/{slug}/qr.png", web.qrHandler).Methods(http.MethodGet)
	r.Handle("/{slug}/edit", web.authenticate(web.editHandler)).Methods(http.MethodGet)
	r.HandleFunc("/{slug}", web.redirectHandler).Methods(http.MethodGet)
	r.HandleFunc("/{slug}", web.authenticate(web.updateHandler)).Methods(http.MethodPost, http.MethodPut, http.MethodPatch)
	r.HandleFunc("/oauth/google/callback", web.oauthCallbackHandler).Methods(http.MethodGet)
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css/"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("js/"))))
	http.Handle("/", r)
	return http.ListenAndServe(port, nil)
}
