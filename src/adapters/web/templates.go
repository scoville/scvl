package web

import (
	"html/template"
	"net/http"
)

var templates = map[string]*template.Template{}
var baseTplPath = "templates"

func renderTemplate(w http.ResponseWriter, r *http.Request, path string, data map[string]interface{}) {
	data["Digest"] = Digest
	tpl := findTemplate("/layouts.tpl", path)
	tpl.ExecuteTemplate(w, "base", data)
}

// cache templates so that it doesn't parse files every time in production
func findTemplate(basePath, path string) (tpl *template.Template) {
	if Digest == "" {
		tpl = template.Must(template.ParseFiles(baseTplPath+basePath, baseTplPath+path))
		return
	}
	tpl, ok := templates[path]
	if ok {
		return
	}
	tpl = template.Must(template.ParseFiles(baseTplPath+basePath, baseTplPath+path))
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
