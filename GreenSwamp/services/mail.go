package services

import (
	"net/smtp"
	"strings"
)

func SendEmail(to, subject, body string) error {
	server := env("SMTP_SERVER", "sandbox.smtp.mailtrap.io")
	port := env("SMTP_PORT", "2525")
	user := env("SMTP_USER", "user")
	pass := env("SMTP_PASS", "pass")

	from := env("SENDER_EMAIL", "test@mailtrap.io")
	name := env("SENDER_NAME", "Greenswamp")

	auth := smtp.PlainAuth("", user, pass, server)

	msg := build(from, name, to, subject, body)

	return smtp.SendMail(server+":"+port, auth, from, []string{to}, []byte(msg))
}

func build(from, name, to, subject, body string) string {
	headers := []string{
		"From: " + name + " <" + from + ">",
		"To: " + to,
		"Subject: " + subject,
		"",
	}
	return strings.Join(headers, "\r\n") + "\r\n" + body
}