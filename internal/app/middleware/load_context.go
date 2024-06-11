package middleware

import (
	"context"
	"net/http"
	"ultigamecast/internal/ctxvar"
	"ultigamecast/internal/models"
	"ultigamecast/internal/pathvar"

	"github.com/google/uuid"
	"github.com/justinas/alice"
)

type Services struct {
	Team interface {
		GetTeam(context.Context, string) (*models.Team, error)
	}
	Tournament interface {
		GetTournament(context.Context, string) (*models.Tournament, error)
	}
	Player interface {
		GetPlayer(context.Context, string) (*models.Player, error)
	}
	Game interface {
		GetGame(context.Context, string) (*models.Game, error)
	}
}

func LoadContext(services Services) alice.Constructor {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = context.WithValue(ctx, ctxvar.Path, r.URL.Path)
			ctx = context.WithValue(ctx, ctxvar.HttpMethod, r.Method)
			u, _ := uuid.NewRandom()
			ctx = context.WithValue(ctx, ctxvar.ReqId, u.String())

			ctx = loadValue(ctx, pathvar.TeamSlug(r), ctxvar.Team, services.Team.GetTeam)
			ctx = loadValue(ctx, pathvar.TournamentSlug(r), ctxvar.Tournament, services.Tournament.GetTournament)
			ctx = loadValue(ctx, pathvar.PlayerSlug(r), ctxvar.Player, services.Player.GetPlayer)
			ctx = loadValue(ctx, pathvar.GameSlug(r), ctxvar.Game, services.Game.GetGame)

			*r = *r.WithContext(ctx)
			h.ServeHTTP(w, r)
		})
	}
}

func loadValue[T any](ctx context.Context, pathVal string, cv ctxvar.ContextVar, getFn func(context.Context, string) (*T, error)) context.Context {
	if pathVal == "" {
		return ctx
	}

	if m, err := getFn(ctx, pathVal); err != nil || m == nil {
		return ctx
	} else {
		return context.WithValue(ctx, cv, m)
	}
}
