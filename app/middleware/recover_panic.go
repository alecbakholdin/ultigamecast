package middleware

import (
	"log/slog"
	"net/http"
)

func RecoverPanic(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				slog.ErrorContext(r.Context(), "Recovered panicking goroutine", "r", rec)
			}
		}()
		h.ServeHTTP(w, r)
	})
}
