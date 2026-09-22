package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/mytheresa/go-hiring-challenge/app/ports/custom_error"
)

type Error struct {
	Message string `json:"error"`
}

func ErrorResponseFromError(w http.ResponseWriter, err error) {
	ErrorResponse(w, StatusCode(err), err.Error())
}

func StatusCode(err error) int {
	switch {
	case errors.Is(err, custom_error.ErrInvalidProductLimit), errors.Is(err, custom_error.ErrCategoryInputInvalid):
		return http.StatusBadRequest
	case errors.Is(err, custom_error.ErrProductNotFound), errors.Is(err, custom_error.ErrCategoryNotFound):
		return http.StatusNotFound
	case errors.Is(err, custom_error.ErrCategoryCodeAlreadyExists):
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
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
