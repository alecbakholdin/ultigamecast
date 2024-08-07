package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"ultigamecast/internal/ctxvar"
	"ultigamecast/internal/models"

	"github.com/justinas/alice"
)

type AuthService interface {
	VerifyJwt(jwt string) (*models.User, error)
}

func LoadUser(a AuthService) alice.Constructor {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			accessToken, err := r.Cookie("access_token")
			if err == nil {
				user, err := a.VerifyJwt(accessToken.Value)
				if err != nil {
					slog.ErrorContext(r.Context(), "invalid jwt", "err", err)
				} else {
					*r = *r.WithContext(context.WithValue(r.Context(), ctxvar.User, user))
				}
			}
			h.ServeHTTP(w, r)
		})
	}
}
