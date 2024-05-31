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

type DatumService interface {
	GetDatum(ctx context.Context, id int64) (*models.TournamentDatum, error)
}

func LoadDatum(t DatumService) alice.Constructor {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			datumSlug := pathvar.TournamentDatumSlug(r)
			if datumSlug == "" {
				http.Error(w, "missing datum identifier", http.StatusBadRequest)
			} else if datum, err := t.GetDatum(r.Context(), 0); errors.Is(service.ErrNotFound, err) {
				http.NotFound(w, r)
			} else if err != nil {
				http.Error(w, "unexpected error", http.StatusInternalServerError)
			} else {
				ctx := context.WithValue(r.Context(), ctxvar.TournamentDatum, datum)
				*r = *r.WithContext(ctx)
			}
			h.ServeHTTP(w, r)
		})
	}
}
