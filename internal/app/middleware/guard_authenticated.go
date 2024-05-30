package middleware

import (
	"net/http"
	"ultigamecast/internal/ctxvar"
)

func GuardAuthenticated(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !ctxvar.IsAuthenticated(r.Context()) {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		h.ServeHTTP(w, r)
	})
}
