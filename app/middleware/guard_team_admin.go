package middleware

import (
	"net/http"
	"ultigamecast/app/ctxvar"
)

func GuardTeamAdmin(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := ctxvar.GetUser(r.Context())
		team := ctxvar.GetTeam(r.Context())
		if team == nil {
			http.NotFound(w, r)
		} else if team.Owner != user.ID {
			http.Error(w, "You are not a team admin", http.StatusForbidden)
		}
		h.ServeHTTP(w, r)
	})
}
