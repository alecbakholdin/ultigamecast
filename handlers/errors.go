package handlers

import (
	"errors"
	"log/slog"
	"net/http"
)

var (
	ErrorInternalServerError = errors.New("unexpected error")
)

func unexpectedError(w http.ResponseWriter, r *http.Request, err error) {
	slog.ErrorContext(r.Context(), "Unexpected error", "err", err)
	http.Error(w, "Unexpected Error", http.StatusInternalServerError)
}
