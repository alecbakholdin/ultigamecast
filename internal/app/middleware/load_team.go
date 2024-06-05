package middleware

import (
	"context"
	"errors"
	"net/http"
	"ultigamecast/internal/app/service"
	"ultigamecast/internal/ctxvar"
	"ultigamecast/internal/models"
	"ultigamecast/internal/pathvar"

	"github.com/justinas/alice"
)

type TeamService interface {
	GetTeam(ctx context.Context, slug string) (*models.Team, error)
	IsTeamAdmin(ctx context.Context) bool
}

func LoadTeam(t TeamService) alice.Constructor {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			teamSlug := pathvar.TeamSlug(r)
			if teamSlug == "" {
				http.Error(w, "missing team identifier", http.StatusBadRequest)
			} else if team, err := t.GetTeam(r.Context(), teamSlug); errors.Is(err, service.ErrNotFound) {
				http.NotFound(w, r)
			} else if err != nil {
				http.Error(w, "unexpected error", http.StatusInternalServerError)
			} else {
				ctx := context.WithValue(r.Context(), ctxvar.Team, team)
				ctx = context.WithValue(ctx, ctxvar.Admin, t.IsTeamAdmin(ctx))
				*r = *r.WithContext(ctx)
			}
			h.ServeHTTP(w, r)
		})
	}
}
