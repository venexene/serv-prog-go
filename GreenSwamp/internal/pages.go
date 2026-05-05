package internal

import (
	"html/template"
	"net/http"
	"strings"
	"time"

	csrf "github.com/gorilla/csrf"
)

type Page struct {
	Title  string
	Active string
	Year   int
}

func (a *App) page(w http.ResponseWriter, r *http.Request) {
	name := cleanPath(r.URL.Path)

	tpl, ok := a.templates[name]
	if !ok {
		name = "404"
		tpl = a.templates["404"]
		w.WriteHeader(http.StatusNotFound)
	}

	data := Page{
		Title:  makeTitle(name),
		Active: name,
		Year:   time.Now().Year(),
	}

	a.render(w, r, tpl, name, data)
}

func (a *App) render(w http.ResponseWriter, r *http.Request, t *template.Template, name string, data Page) {
	err := t.ExecuteTemplate(w, name+".html", map[string]any{
		"Page": data,
		"CSRF": template.HTML(csrf.TemplateField(r)),
	})
	if err != nil {
		http.Error(w, "template error", 500)
	}
}

func cleanPath(p string) string {
	p = strings.Trim(p, "/")
	if p == "" {
		return "index"
	}
	if p == "terms" {
		return "tos"
	}
	return p
}

func makeTitle(p string) string {
	switch p {
	case "index":
		return "Greenswamp - Home"
	case "about":
		return "Greenswamp - About"
	case "contact":
		return "Greenswamp - Contact"
	case "privacy":
		return "Greenswamp - Privacy"
	case "tos":
		return "Greenswamp - Terms"
	case "404":
		return "Not found"
	default:
		return "Greenswamp"
	}
}