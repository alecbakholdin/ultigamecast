package service

import (
	"errors"
	"testing"
	"ultigamecast/test/testdb"

	"github.com/stretchr/testify/assert"
)

func TestSignUp(t *testing.T) {
	q, _ := testdb.DB()
	a := NewAuth(q, "secret")
	jwt, err := a.SignUp("email@email.com", "randompassword")
	if err != nil {
		t.Fatalf("could not sign up: %s", err)
	}
	user, err := a.VerifyJwt(jwt)
	if err != nil{
		t.Fatalf("could not verify jwt: %s", err)
	}
	assert.Equal(t, "email@email.com", user.Email)
}

func TestSignInWithPasswordOnlyWorksWithProperPassword(t *testing.T) {
	q, _ := testdb.DB()
	a := NewAuth(q, "secret")
	_, err := a.SignUp("email2@email.com", "randompassword")
	if err != nil {
		t.Fatalf("could not sign up: %s", err)
	}
	jwt, err := a.SignInWithPassword("email2@email.com", "randompassword")
	if err != nil {
		t.Fatalf("could not log in: %s", err)
	}
	user, err := a.VerifyJwt(jwt)
	if err != nil {
		t.Fatalf("could not verify jwt: %s", err)
	}
	assert.Equal(t, "email2@email.com", user.Email)
	
	_, err = a.SignInWithPassword("email2@email.com", "wrongpassword")
	assert.Truef(t, errors.Is(err, ErrInvalidCredentials), "login should have failed since password is wrong, but didn't. Error: %s", err)
}