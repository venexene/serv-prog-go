package services

import (
	"encoding/csv"
	"os"
	"path/filepath"
	"time"
)

func SaveContact(name, email, topic, message string) error {
	return writeCSV(
		env("CONTACTS_CSV", "data/contacts.csv"),
		[]string{"timestamp", "name", "email", "topic", "message"},
		[]string{now(), name, email, topic, message},
	)
}

func SaveSubscriber(email string) error {
	return writeCSV(
		env("SUBSCRIBERS_CSV", "data/subscribers.csv"),
		[]string{"timestamp", "email"},
		[]string{now(), email},
	)
}

func writeCSV(path string, header, row []string) error {
	os.MkdirAll(filepath.Dir(path), 0755)

	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	info, _ := f.Stat()
	if info.Size() == 0 {
		w.Write(header)
	}

	return w.Write(row)
}

func now() string {
	return time.Now().Format(time.RFC3339)
}

func env(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}