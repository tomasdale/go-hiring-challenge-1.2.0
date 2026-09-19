package api

import (
	"encoding/json"
	"net/http"
)

type Error struct {
	Message string `json:"error"`
}

func OKResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "Failed to encode response")
		return
	}
}

func ErrorResponse(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(Error{Message: message})
}
