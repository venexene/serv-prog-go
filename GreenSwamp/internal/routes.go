package internal

import "net/http"

func (a *App) Routes(mux *http.ServeMux) {
	mux.HandleFunc("/", a.page)
	mux.HandleFunc("/api/contact", a.contact)
	mux.HandleFunc("/api/subscribe", a.subscribe)

	if a.config.StaticDir != "" {
		fs := http.FileServer(http.Dir(a.config.StaticDir))
		mux.Handle("/static/", http.StripPrefix("/static/", fs))
	}
}