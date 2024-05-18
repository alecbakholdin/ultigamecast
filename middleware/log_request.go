package middleware

import (
	"log/slog"
	"net/http"
)

func LogRequest(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.InfoContext(r.Context(), "Request")
		h.ServeHTTP(w, r)
	})
}
