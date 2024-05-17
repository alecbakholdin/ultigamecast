package handlers

import (
	"net/http"
	"ultigamecast/models"
	view_auth "ultigamecast/view/auth"
)

type Auth struct {
	authService AuthService
}

type AuthService interface {
	SignIn(email, password string) (*models.User, error)
	SignOut(*models.User) 
}

func NewAuth(a AuthService) *Auth {
	return &Auth{
		authService: a,
	}
}

func (a *Auth) GetLogin(w http.ResponseWriter, r *http.Request) {
	//var validate = validator.New(validator.WithRequiredStructEnabled())
	view_auth.LoginPage().Render(r.Context(), w);
}

func (a *Auth) PostLogin(w http.ResponseWriter, r *http.Request) {

}