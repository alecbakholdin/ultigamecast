package middleware

import (
	"context"
	"net/http"
	"ultigamecast/internal/ctxvar"

	"github.com/google/uuid"
	"github.com/justinas/alice"
)



func LoadContext(t TeamService) alice.Constructor {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = context.WithValue(ctx, ctxvar.Path, r.URL.Path)
			ctx = context.WithValue(ctx, ctxvar.HttpMethod, r.Method)
			u, _ := uuid.NewRandom()
			ctx = context.WithValue(ctx, ctxvar.ReqId, u.String())
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
