package middleware

import (
	"net/http"
	"ultigamecast/app/ctxvar"
)

func GuardTeamAdmin(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !ctxvar.IsAdmin(r.Context()) {
			http.Error(w, "You are not a team admin", http.StatusForbidden)
			return
		}
		h.ServeHTTP(w, r)
	})
}
