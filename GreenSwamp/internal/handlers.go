package internal

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	svc "github.com/venexene/serv-prog-go/greenswamp/services"
)

func (a *App) subscribe(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		jsonErr(w, "method not allowed", 405)
		return
	}

	email := strings.TrimSpace(r.FormValue("email"))

	if !regexp.MustCompile(`.+@.+\..+`).MatchString(email) {
		jsonErr(w, "invalid email", 400)
		return
	}

	if err := svc.SaveSubscriber(email); err != nil {
		jsonErr(w, "save error", 500)
		return
	}

	go svc.SendEmail(email, "Welcome!", "Stay Froggy!")

	jsonOK(w, "subscribed")
}

func (a *App) contact(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		jsonErr(w, "method not allowed", 405)
		return
	}

	name := strings.TrimSpace(r.FormValue("name"))
	email := strings.TrimSpace(r.FormValue("email"))
	topic := strings.TrimSpace(r.FormValue("topic"))
	msg := strings.TrimSpace(r.FormValue("message"))

	if name == "" || email == "" || msg == "" {
		jsonErr(w, "missing fields", 400)
		return
	}

	if !strings.HasSuffix(strings.ToLower(email), ".edu") {
		jsonErr(w, "must be .edu email", 400)
		return
	}

	if err := svc.SaveContact(name, email, topic, msg); err != nil {
		jsonErr(w, "save error", 500)
		return
	}

	go func() {
		admin := getEnv("ADMIN_EMAIL", "admin@greenswamp.com")
		body := fmt.Sprintf("From: %s (%s)\n\n%s", name, email, msg)
		_ = svc.SendEmail(admin, "New contact", body)
	}()

	jsonOK(w, "received")
}