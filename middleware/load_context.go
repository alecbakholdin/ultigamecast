package middleware

import (
	"context"
	"net/http"
	"ultigamecast/app/ctx_var"
	"ultigamecast/models"

	"github.com/google/uuid"
	"github.com/justinas/alice"
)

type TeamService interface {
	GetBySlug(ctx context.Context, slug string) (*models.Team, error)
}

type TournamentService interface {
	GetBySlug(ctx context.Context, team *models.Team, slug string) (*models.Tournament)
}

type GameService interface {
	GetById(ctx context.Context, tournament *models.Tournament, id string) (*models.Game)
}

func LoadContext(t TeamService) alice.Constructor{
	return func (h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = context.WithValue(ctx, ctx_var.Path, r.URL.Path)
			ctx = context.WithValue(ctx, ctx_var.HttpMethod, r.Method)
			u, _ := uuid.NewRandom()
			ctx = context.WithValue(ctx, ctx_var.ReqId, u.String())
			if teamSlug := r.PathValue("teamSlug"); teamSlug != "" {
				if team, err := t.GetBySlug(ctx, teamSlug); err == nil{
					ctx = context.WithValue(ctx, ctx_var.Team, team)
				}
			}
			// if tournamentSlug := r.PathValue("tournamentSlug"); tournamentSlug != "" {
			// 	if tournament, err := t.GetBySlug(tournamentSlug); err == nil{
			// 		ctx = context.WithValue(ctx, ctx_var.Tournament, tournament)
			// 	}
			// }
			// if gameSlug := r.PathValue("gameSlug"); gameSlug != "" {
			// 	if game, err := t.GetBySlug(gameSlug); err == nil{
			// 		ctx = context.WithValue(ctx, ctx_var.Game, game)
			// 	}
			// }
			*r = *r.WithContext(ctx)
			h.ServeHTTP(w, r)
		})
	}
	}