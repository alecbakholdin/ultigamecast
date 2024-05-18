package middleware

import (
	"context"
	"net/http"
	"ultigamecast/app/ctx_var"

	"github.com/google/uuid"
)

func LoadContext(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = context.WithValue(ctx, ctx_var.Path, r.URL.Path)
		ctx = context.WithValue(ctx, ctx_var.HttpMethod, r.Method)
		u, _ := uuid.NewRandom()
		ctx = context.WithValue(ctx, ctx_var.ReqId, u.String())
		*r = *r.WithContext(ctx)
		h.ServeHTTP(w, r)
	})
}
