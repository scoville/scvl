package web

import (
	"net/http"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/scoville/scvl/src/domain"
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
func New(e *engine.Engine, sessionSecret, mainDomain string) *Web {
	w := &Web{
		engine: e,
		store:  sessions.NewCookieStore([]byte(sessionSecret)),
	}
	w.store.Options.Domain = mainDomain
	return w
}

// Start starts listen and serve
func (web *Web) Start(port string) error {
	r := mux.NewRouter()

	r.Handle("/", web.authenticate(web.rootHandler)).Methods(http.MethodGet)
	r.Handle("/logout", web.authenticate(web.logoutHandler)).Methods(http.MethodPost)
	r.Handle("/shorten", web.authenticate(web.shortenHandler)).Methods(http.MethodPost)
	r.HandleFunc("/api/shorten", web.shortenByAPIHandler).Methods(http.MethodPost)
	r.Handle("/api/key", web.authenticate(web.publishAPIKeyHandler)).Methods(http.MethodPost)
	r.Handle("/pages", web.authenticate(web.pagesHandler)).Methods(http.MethodGet)
	r.Handle("/files", web.authenticate(web.filesHandler)).Methods(http.MethodGet)
	r.Handle("/files", web.authenticate(web.fileUploadHandler)).Methods(http.MethodPost)
	r.HandleFunc("/files/{slug}", web.fileShowHandler).Methods(http.MethodGet)
	r.Handle("/files/{slug}/download", web.authenticate(web.fileDownloadHandler)).Methods(http.MethodPost)
	r.Handle("/files/{slug}/edit", web.authenticate(web.editFileHandler)).Methods(http.MethodGet)
	r.HandleFunc("/files/{slug}", web.authenticate(web.updateFileHandler)).Methods(http.MethodPost, http.MethodPut, http.MethodPatch)
	r.Handle("/images", web.authenticate(web.imagesHandler)).Methods(http.MethodGet)
	r.Handle("/images", web.authenticate(web.imageUploadHandler)).Methods(http.MethodPost)

	r.Handle("/emails", web.authenticate(web.emailsHandler)).Methods(http.MethodGet)
	r.Handle("/emails", web.authenticate(web.emailCreateHandler)).Methods(http.MethodPost)
	r.Handle("/emails/send", web.authenticate(web.emailSendHandler)).Methods(http.MethodPost)
	r.HandleFunc("/emails/{id:[0-9]+}/read", web.emailReadHandler).Methods(http.MethodGet)

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

func (web *Web) rootHandler(w http.ResponseWriter, r *http.Request) {
	resp := map[string]interface{}{}
	user, ok := context.Get(r, "user").(*domain.User)
	if ok {
		resp["User"] = user
	}
	loginURL, ok := context.Get(r, "login_url").(string)
	if ok {
		resp["LoginURL"] = loginURL
	}
	renderTemplate(w, r, "/index.tpl", resp)
}

func (web *Web) logoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := web.store.Get(r, "scvl")
	delete(session.Values, "user_id")
	web.store.Save(r, w, session)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
