package internal

import (
	"encoding/json"
	"net/http"
)

func jsonOK(w http.ResponseWriter, msg string) {
	writeJSON(w, 200, map[string]string{
		"status":  "ok",
		"message": msg,
	})
}

func jsonErr(w http.ResponseWriter, msg string, code int) {
	writeJSON(w, code, map[string]string{
		"status":  "error",
		"message": msg,
	})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}