package web

import (
	"html/template"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/context"
	"github.com/scoville/scvl/src/domain"
)

var templates = map[string]*template.Template{}
var baseTplPath = "templates"

var tempFuncs = map[string]interface{}{
	// Replaces newlines with <br>
	"nl2br": func(text string) template.HTML {
			return template.HTML(strings.Replace(template.HTMLEscapeString(text), "\n", "<br>", -1))
	},
	// Skips sanitation on the parameter.  Do not use with dynamic data.
	"raw": func(text string) template.HTML {
			return template.HTML(text)
	},
}

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
	data["IsMainHost"] = (r.Host == os.Getenv("MAIN_DOMAIN"))
	data["IsFileHost"] = (r.Host == os.Getenv("FILE_DOMAIN"))
	data["IsImageHost"] = (r.Host == os.Getenv("IMAGE_DOMAIN"))
	data["IsEmailHost"] = (r.Host == os.Getenv("EMAIL_DOMAIN"))
	data["MainHost"] = scheme + os.Getenv("MAIN_DOMAIN")
	data["FileHost"] = scheme + os.Getenv("FILE_DOMAIN")
	data["ImageHost"] = scheme + os.Getenv("IMAGE_DOMAIN")
	data["EmailHost"] = scheme + os.Getenv("EMAIL_DOMAIN")
	tpl := findTemplate("/layouts.html", path)
	tpl.ExecuteTemplate(w, "base", data)
}

// cache templates so that it doesn't parse files every time in production
func findTemplate(basePath, path string) (tpl *template.Template) {
	if Digest == "" {
		tpl = template.Must(template.New(path).Funcs(tempFuncs).ParseFiles(baseTplPath+basePath, baseTplPath+path))
		return
	}
	tpl, ok := templates[path]
	if ok {
		return
	}
	tpl = template.Must(template.New(path).Funcs(tempFuncs).ParseFiles(baseTplPath+basePath, baseTplPath+path))
	templates[path] = tpl
	return
}

func findTemplateWithoutBase(path string) (tpl *template.Template, err error) {
	paths := strings.Split(path, "/")
	filename := paths[len(paths) - 1]
	if Digest == "" {
		tpl, err = template.New(filename).Funcs(tempFuncs).ParseFiles(baseTplPath + path)
		return
	}
	tpl, ok := templates[path]
	if ok {
		return
	}
	tpl, err = template.New(filename).Funcs(tempFuncs).ParseFiles(baseTplPath + path)
	templates[path] = tpl
	return
}
