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

type TournamentService interface {
	GetTournament(ctx context.Context, slug string) (*models.Tournament, error)
}

func LoadTournament(t TournamentService) alice.Constructor {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tournamentSlug := pathvar.TournamentSlug(r)
			if tournamentSlug == "" {
				http.Error(w, "missing tournament identifier", http.StatusBadRequest)
			} else if tournament, err := t.GetTournament(r.Context(), tournamentSlug); errors.Is(service.ErrNotFound, err) {
				http.NotFound(w, r)
			} else if err != nil {
				http.Error(w, "unexpected error", http.StatusInternalServerError)
			} else {
				ctx := context.WithValue(r.Context(), ctxvar.Tournament, tournament)
				*r = *r.WithContext(ctx)
			}
			h.ServeHTTP(w, r)
		})
	}
}
