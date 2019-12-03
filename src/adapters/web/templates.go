package web

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/context"
	"github.com/scoville/scvl/src/domain"
)

var templates = map[string]*template.Template{}
var baseTplPath = "templates"

func renderTemplate(w http.ResponseWriter, r *http.Request, path string, data map[string]interface{}) {
	data["Digest"] = Digest
	scheme := "https://"
	if strings.Contains(os.Getenv("MAIN_DOMAIN"), "localhost") {
		scheme = "http://"
	}
	user, ok := context.Get(r, "user").(*domain.User)
	if ok {
		data["User"] = user
		data["UserSignedIn"] = true
	}
	log.Println(r.Host)
	data["IsMainHost"] = (r.Host == os.Getenv("MAIN_DOMAIN"))
	data["IsFileHost"] = (r.Host == os.Getenv("FILE_DOMAIN"))
	data["IsImageHost"] = (r.Host == os.Getenv("IMAGE_DOMAIN"))
	data["MainHost"] = scheme + os.Getenv("MAIN_DOMAIN")
	data["FileHost"] = scheme + os.Getenv("FILE_DOMAIN")
	data["ImageHost"] = scheme + os.Getenv("IMAGE_DOMAIN")
	tpl := findTemplate("/layouts.tpl", "/login.tpl", path)
	tpl.ExecuteTemplate(w, "base", data)
}

// cache templates so that it doesn't parse files every time in production
func findTemplate(basePath, loginPath, path string) (tpl *template.Template) {
	if Digest == "" {
		tpl = template.Must(template.ParseFiles(baseTplPath+basePath, baseTplPath+loginPath, baseTplPath+path))
		return
	}
	tpl, ok := templates[path]
	if ok {
		return
	}
	tpl = template.Must(template.ParseFiles(baseTplPath+basePath, baseTplPath+loginPath, baseTplPath+path))
	templates[path] = tpl
	return
}

func findTemplateWithoutBase(path string) (tpl *template.Template) {
	if Digest == "" {
		tpl = template.Must(template.ParseFiles(baseTplPath + path))
		return
	}
	tpl, ok := templates[path]
	if ok {
		return
	}
	tpl = template.Must(template.ParseFiles(baseTplPath + path))
	templates[path] = tpl
	return
}
