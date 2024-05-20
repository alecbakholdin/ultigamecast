package middleware

import (
	"net/http"
	"ultigamecast/app/ctx_var"
)

func GuardAuthenticated(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !ctx_var.IsAuthenticated(r.Context()) {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		h.ServeHTTP(w, r)
	})
}
