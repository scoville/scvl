package main

import (
	"html/template"
	"net/http"
)

var templates = map[string]*template.Template{}

func renderTemplate(w http.ResponseWriter, r *http.Request, path string, data map[string]interface{}) {
	tpl := findTemplate("/layouts.tpl", path)
	tpl.ExecuteTemplate(w, "base", data)
}

// cache templates so that it doesn't parse files every time in production
func findTemplate(basePath, path string) (tpl *template.Template) {
	if config.env == envDev {
		tpl = template.Must(template.ParseFiles(config.baseTplPath+basePath, config.baseTplPath+path))
		return
	}
	tpl, ok := templates[path]
	if ok {
		return
	}
	tpl = template.Must(template.ParseFiles(config.baseTplPath+basePath, config.baseTplPath+path))
	templates[path] = tpl
	return
}
