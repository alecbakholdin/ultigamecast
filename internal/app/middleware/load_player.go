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

type PlayerService interface {
	GetPlayer(ctx context.Context, slug string) (*models.Player, error)
}

func LoadPlayer(t PlayerService) alice.Constructor {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			playerSlug := pathvar.PlayerSlug(r)
			if playerSlug == "" {
				http.Error(w, "missing player identifier", http.StatusBadRequest)
			} else if player, err := t.GetPlayer(r.Context(), playerSlug); errors.Is(service.ErrNotFound, err) {
				http.NotFound(w, r)
			} else if err != nil {
				http.Error(w, "unexpected error", http.StatusInternalServerError)
			} else {
				ctx := context.WithValue(r.Context(), ctxvar.Player, player)
				*r = *r.WithContext(ctx)
			}
			h.ServeHTTP(w, r)
		})
	}
}
