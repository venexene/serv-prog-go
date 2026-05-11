package internal

import (
	"html/template"
	"log"
	"net/http"

	csrf "github.com/gorilla/csrf"
)

type App struct {
	log       *log.Logger
	templates map[string]*template.Template
	config    Config
}

func CreateApp(cfg Config, logger *log.Logger) *App {
	tpls, err := LoadTemplates(cfg.TplDir)
	if err != nil {
		logger.Fatal(err)
	}

	return &App{
		log:       logger,
		templates: tpls,
		config:    cfg,
	}
}

func (a *App) WithMiddleware(next http.Handler) http.Handler {
	csrfKey := getEnv("CSRF_KEY", "very-secret-key-123456789012")

	return csrf.Protect(
		[]byte(csrfKey),
		csrf.Secure(false),
		csrf.TrustedOrigins([]string{
			"localhost:8080",
			"127.0.0.1:8080",
		}),
	)(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			a.log.Printf("%s %s", r.Method, r.URL.Path)
			next.ServeHTTP(w, r)
		}),
	)
}