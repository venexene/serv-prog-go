package internal

import (
	"html/template"
	"path/filepath"
)

func LoadTemplates(dir string) (map[string]*template.Template, error) {
	base := template.New("base")

	shared, _ := filepath.Glob(filepath.Join(dir, "shared/*.html"))
	if len(shared) > 0 {
		base.ParseFiles(shared...)
	}

	pages := []string{"index", "about", "contact", "privacy", "tos", "404"}

	result := make(map[string]*template.Template)

	for _, p := range pages {
		t, _ := base.Clone()
		t.ParseFiles(filepath.Join(dir, "pages", p+".html"))
		result[p] = t
	}

	return result, nil
}