package account

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	csrf "github.com/gorilla/csrf"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AccountController struct {
	repo            *Repository
	db              *gorm.DB
	tmpl            *template.Template
	logger          *log.Logger
	cookieName      string
	defaultRedirect string
}

type AccountPageData struct {
	Title       string
	Error       string
	Next        string
	Login       string
	Username    string
	DisplayName string
	Email       string
	CSRFField   template.HTML
}

func RegisterRoutes(mux *http.ServeMux, db *gorm.DB, logger *log.Logger, templatesDir string) {
	ctrl := &AccountController{
		repo:            NewRepository(db),
		db:              db,
		tmpl:            mustTemplates(templatesDir),
		logger:          logger,
		cookieName:      "gs_auth_session",
		defaultRedirect: "/posts/feed",
	}

	mux.HandleFunc("/login", ctrl.login)
	mux.HandleFunc("/register", ctrl.register)
	mux.HandleFunc("/logout", ctrl.logout)
}

func mustTemplates(dir string) *template.Template {
	funcMap := template.FuncMap{}
	pattern := dir + "/*.html"
	tmpl, err := template.New("account").Funcs(funcMap).ParseGlob(pattern)
	if err != nil {
		log.Fatalf("failed to parse account templates: %v", err)
	}
	return tmpl
}

func (c *AccountController) login(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		c.render(w, r, "login.html", AccountPageData{
			Title:     "Login",
			Next:      safeNext(r.URL.Query().Get("next"), c.defaultRedirect),
			CSRFField: template.HTML(csrf.TemplateField(r)),
		})
	case http.MethodPost:
		_ = r.ParseForm()

		login := strings.TrimSpace(r.FormValue("login"))
		password := r.FormValue("password")
		next := safeNext(r.FormValue("next"), c.defaultRedirect)

		if login == "" || password == "" {
			c.render(w, r, "login.html", AccountPageData{
				Title:     "Login",
				Error:     "Please enter your login and password.",
				Next:      next,
				Login:     login,
				CSRFField: template.HTML(csrf.TemplateField(r)),
			})
			return
		}

		user, err := c.repo.FindByLogin(r.Context(), login)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.render(w, r, "login.html", AccountPageData{
					Title:     "Login",
					Error:     "User not found.",
					Next:      next,
					Login:     login,
					CSRFField: template.HTML(csrf.TemplateField(r)),
				})
				return
			}
			http.Error(w, "database error", http.StatusInternalServerError)
			return
		}

		if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
			c.render(w, r, "login.html", AccountPageData{
				Title:     "Login",
				Error:     "Incorrect password.",
				Next:      next,
				Login:     login,
				CSRFField: template.HTML(csrf.TemplateField(r)),
			})
			return
		}

		token, err := newToken()
		if err != nil {
			http.Error(w, "token error", http.StatusInternalServerError)
			return
		}

		expiresAt := time.Now().Add(7 * 24 * time.Hour)

		if err := c.repo.CreateSession(r.Context(), user.UserID, token, expiresAt); err != nil {
			http.Error(w, "session error", http.StatusInternalServerError)
			return
		}
		_ = c.repo.UpdateLastLogin(r.Context(), user.UserID)

		c.setCookie(w, r, token, expiresAt)
		http.Redirect(w, r, next, http.StatusSeeOther)

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (c *AccountController) register(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		c.render(w, r, "register.html", AccountPageData{
			Title:     "Register",
			Next:      safeNext(r.URL.Query().Get("next"), c.defaultRedirect),
			CSRFField: template.HTML(csrf.TemplateField(r)),
		})
	case http.MethodPost:
		_ = r.ParseForm()

		username := strings.TrimSpace(r.FormValue("username"))
		displayName := strings.TrimSpace(r.FormValue("display_name"))
		email := strings.TrimSpace(r.FormValue("email"))
		password := r.FormValue("password")
		password2 := r.FormValue("password_confirm")
		next := safeNext(r.FormValue("next"), c.defaultRedirect)

		if username == "" || email == "" || password == "" || password2 == "" {
			c.render(w, r, "register.html", AccountPageData{
				Title:       "Register",
				Error:       "Please fill in all required fields.",
				Next:        next,
				Username:    username,
				DisplayName: displayName,
				Email:       email,
				CSRFField:   template.HTML(csrf.TemplateField(r)),
			})
			return
		}

		if password != password2 {
			c.render(w, r, "register.html", AccountPageData{
				Title:       "Register",
				Error:       "Passwords do not match.",
				Next:        next,
				Username:    username,
				DisplayName: displayName,
				Email:       email,
				CSRFField:   template.HTML(csrf.TemplateField(r)),
			})
			return
		}

		if len(password) < 8 {
			c.render(w, r, "register.html", AccountPageData{
				Title:       "Register",
				Error:       "Password must be at least 8 characters.",
				Next:        next,
				Username:    username,
				DisplayName: displayName,
				Email:       email,
				CSRFField:   template.HTML(csrf.TemplateField(r)),
			})
			return
		}

		if displayName == "" {
			displayName = username
		}

		if _, err := c.repo.FindByUsername(r.Context(), username); err == nil {
			c.render(w, r, "register.html", AccountPageData{
				Title:       "Register",
				Error:       "That username is already taken.",
				Next:        next,
				Username:    username,
				DisplayName: displayName,
				Email:       email,
				CSRFField:   template.HTML(csrf.TemplateField(r)),
			})
			return
		}
		if _, err := c.repo.FindByEmail(r.Context(), email); err == nil {
			c.render(w, r, "register.html", AccountPageData{
				Title:       "Register",
				Error:       "That email is already registered.",
				Next:        next,
				Username:    username,
				DisplayName: displayName,
				Email:       email,
				CSRFField:   template.HTML(csrf.TemplateField(r)),
			})
			return
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "password error", http.StatusInternalServerError)
			return
		}

		user := &IdentityUser{
			Username:     username,
			DisplayName:  displayName,
			Email:        email,
			PasswordHash: string(hash),
			IsActive:     true,
		}

		if err := c.repo.CreateUser(r.Context(), user); err != nil {
			http.Error(w, "create user error", http.StatusInternalServerError)
			return
		}

		err = c.db.WithContext(r.Context()).Exec(
			`INSERT INTO users (user_id, username, display_name, is_active, created_at) 
			 VALUES (?, ?, ?, ?, ?)`,
			user.UserID,
			user.Username,
			user.DisplayName,
			true,
			time.Now(),
		).Error
		if err != nil {
			c.logger.Printf("warning: failed to create posts.users record: %v", err)
		}

		token, err := newToken()
		if err != nil {
			http.Error(w, "token error", http.StatusInternalServerError)
			return
		}
		expiresAt := time.Now().Add(7 * 24 * time.Hour)

		if err := c.repo.CreateSession(r.Context(), user.UserID, token, expiresAt); err != nil {
			http.Error(w, "session error", http.StatusInternalServerError)
			return
		}
		_ = c.repo.UpdateLastLogin(r.Context(), user.UserID)

		c.setCookie(w, r, token, expiresAt)
		http.Redirect(w, r, next, http.StatusSeeOther)

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (c *AccountController) logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	target := safeNext(r.URL.Query().Get("next"), c.defaultRedirect)
	if target == c.defaultRedirect {
		if ref := sameSiteReferer(r); ref != "" {
			target = ref
		}
	}

	if cookie, err := r.Cookie(c.cookieName); err == nil && cookie.Value != "" {
		_ = c.repo.RevokeSession(r.Context(), cookie.Value)
	}

	c.clearCookie(w, r)
	http.Redirect(w, r, target, http.StatusSeeOther)
}

func (c *AccountController) render(w http.ResponseWriter, r *http.Request, name string, data AccountPageData) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := c.tmpl.ExecuteTemplate(w, name, data); err != nil {
		c.logger.Printf("account template %s: %v", name, err)
		http.Error(w, "template error", http.StatusInternalServerError)
	}
}

func (c *AccountController) setCookie(w http.ResponseWriter, r *http.Request, token string, expiresAt time.Time) {
	http.SetCookie(w, &http.Cookie{
		Name:     c.cookieName,
		Value:    token,
		Path:     "/",
		Expires:  expiresAt,
		MaxAge:   int(time.Until(expiresAt).Seconds()),
		HttpOnly: true,
		Secure:   r.TLS != nil,
		SameSite: http.SameSiteLaxMode,
	})
}

func (c *AccountController) clearCookie(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     c.cookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   r.TLS != nil,
		SameSite: http.SameSiteLaxMode,
	})
}

func newToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func safeNext(next, fallback string) string {
	next = strings.TrimSpace(next)
	if next == "" {
		return fallback
	}
	if strings.HasPrefix(next, "//") {
		return fallback
	}
	if !strings.HasPrefix(next, "/") {
		return fallback
	}
	return next
}

func sameSiteReferer(r *http.Request) string {
	ref := r.Referer()
	if ref == "" {
		return ""
	}

	u, err := url.Parse(ref)
	if err != nil {
		return ""
	}

	if u.Host != "" && !strings.EqualFold(u.Host, r.Host) {
		return ""
	}

	if u.Path == "" {
		return ""
	}
	if u.RawQuery != "" {
		return u.Path + "?" + u.RawQuery
	}
	return u.Path
}