package handlers

import (
	"errors"
	"log/slog"
	"net/http"
	"ultigamecast/internal/app/handlers/htmx"
	"ultigamecast/internal/app/service"
	view_auth "ultigamecast/web/view/auth"
)

type Auth struct {
	authService AuthService
}

type AuthService interface {
	SignInWithPassword(email, password string) (jwt string, err error)
	SignUp(email, password string) (jwt string, err error)
}

func NewAuth(a AuthService) *Auth {
	return &Auth{
		authService: a,
	}
}

func (a *Auth) GetLogin(w http.ResponseWriter, r *http.Request) {
	view_auth.LoginPage().Render(r.Context(), w)
}

func (a *Auth) PostLogin(w http.ResponseWriter, r *http.Request) {
	dto := &view_auth.LoginFormDTO{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}

	if !dto.Validate(dto) {
		view_auth.LoginForm(dto).Render(r.Context(), w)
		return
	}

	if jwt, err := a.authService.SignInWithPassword(dto.Email, dto.Password); errors.Is(err, service.ErrInvalidCredentials) {
		dto.AddFormError("Invalid credential")
		view_auth.LoginForm(dto).Render(r.Context(), w)
	} else if err != nil {
		slog.ErrorContext(r.Context(), "unexpected error authenticating", "err", err)
		dto.AddFormError("Unexpected error")
		view_auth.LoginForm(dto).Render(r.Context(), w)
	} else {
		http.SetCookie(w, &http.Cookie{
			Name:     "access_token",
			Value:    jwt,
			Path:     "/",
			MaxAge:   86400 * 7,
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteLaxMode,
		})
		htmx.HxRedirect(w, "/")
	}
}

func (a *Auth) GetSignup(w http.ResponseWriter, r *http.Request) {
	view_auth.SignUpPage().Render(r.Context(), w)
}

func (a *Auth) PostSignup(w http.ResponseWriter, r *http.Request) {
	dto := &view_auth.SignUpDTO{
		Email:           r.FormValue("email"),
		Password:        r.FormValue("password"),
		ConfirmPassword: r.FormValue("confirm_password"),
	}
	if !dto.Validate(dto) {
		view_auth.SignUpForm(dto).Render(r.Context(), w)
		return
	}

	if jwt, err := a.authService.SignUp(dto.Email, dto.Password); errors.Is(err, service.ErrAccountExists) {
		dto.AddFormError("An account with that email already exists")
		view_auth.SignUpForm(dto).Render(r.Context(), w)
	} else if err != nil {
		slog.ErrorContext(r.Context(), "unexpected error authenticating", "err", err)
		dto.AddFormError("Unexpected error")
		view_auth.SignUpForm(dto).Render(r.Context(), w)
	} else {
		http.SetCookie(w, &http.Cookie{
			Name:     "access_token",
			Value:    jwt,
			Path:     "/",
			MaxAge:   86400 * 7,
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteLaxMode,
		})
		htmx.HxRedirect(w, "/")
	}
}

func (a *Auth) PostLogout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:   "access_token",
		Value:  "deleted",
		MaxAge: 1,
	})
	if r.URL.Path == "/" {
		htmx.HxRefresh(w)
	} else {
		htmx.HxRedirect(w, "/")
	}
}
