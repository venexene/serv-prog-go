package posts

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	account "github.com/venexene/serv-prog-go/greenswamp/account"
	"gorm.io/gorm"
)

type Config struct {
	BasePath     string
	TemplatesDir string
}

func RegisterRoutes(mux *http.ServeMux, db *gorm.DB, logger *log.Logger, cfg Config) {
	if cfg.BasePath == "" {
		cfg.BasePath = "/posts"
	}
	if cfg.TemplatesDir == "" {
		cfg.TemplatesDir = filepath.Join("templates", "posts")
	}

	ctrl := &Controller{
		repo:     NewRepository(db),
		authRepo: account.NewRepository(db),
		tmpl:     mustTemplates(cfg.TemplatesDir),
		basePath: cfg.BasePath,
		logger:   logger,
	}

	mux.HandleFunc(cfg.BasePath+"/feed/post/", ctrl.handlePostDetail)
	mux.HandleFunc(cfg.BasePath+"/profile/", ctrl.handleProfile)
	mux.HandleFunc(cfg.BasePath+"/ponds/", ctrl.handlePond)
	mux.HandleFunc(cfg.BasePath+"/ponds", ctrl.handlePondsIndex)
	mux.HandleFunc(cfg.BasePath+"/api/create", ctrl.handleCreatePost)
	mux.HandleFunc(cfg.BasePath+"/api/interact", ctrl.handleInteract)
	mux.HandleFunc(cfg.BasePath+"/api/comment", ctrl.handleComment)
	mux.HandleFunc(cfg.BasePath+"/feed", ctrl.handleFeed)
	mux.HandleFunc(cfg.BasePath+"/", ctrl.handleIndex)
	mux.HandleFunc(cfg.BasePath, ctrl.handleIndex)
}

func mustTemplates(dir string) *template.Template {
	funcMap := template.FuncMap{
		"avatar":       avatarOrFallback,
		"bio":          bioOrEmpty,
		"mediaKind":    mediaKind,
		"formatTime":   formatTime,
		"uintToString":  uintToString,
	}

	pattern := filepath.Join(dir, "*.html")
	tmpl, err := template.New("posts").Funcs(funcMap).ParseGlob(pattern)
	if err != nil {
		log.Fatalf("failed to parse post templates: %v", err)
	}
	return tmpl
}