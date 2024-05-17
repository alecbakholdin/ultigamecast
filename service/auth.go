package service

import (
	"ultigamecast/models"
)

type Auth struct {
	q *models.Queries
}

func NewAuth(q *models.Queries) *Auth {
	return &Auth{
		q: q,
	}
}

func (a *Auth) SignIn(email string, password string) (*models.User, error) {
	panic("not implemented") // TODO: Implement
}

func (a *Auth) SignOut(_ *models.User) {
	panic("not implemented") // TODO: Implement
}
