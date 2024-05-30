package service

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
	"ultigamecast/internal/models"

	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
	q      *models.Queries
	secret []byte
}

func NewAuth(q *models.Queries, secret string) *Auth {
	return &Auth{
		q:      q,
		secret: []byte(secret),
	}
}

var (
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrFailedToAuth        = errors.New("failed to authenticate")
	ErrAccountExists       = errors.New("email already exists")
	ErrFailedToCreateToken = errors.New("failed to create token")
	ErrExpiredCredentials  = errors.New("expired credentials")
)

// fetches the user and returns the created jwt if successful. Returns
// [ErrInvalidCredentials] if user's email/password doesn't match anything
// or [ErrFailedToAuth] if something unexpected happens
func (a *Auth) SignInWithPassword(email, password string) (jwt string, err error) {
	user, err := a.q.GetUser(context.Background(), email)
	if errors.Is(sql.ErrNoRows, err) {
		return "", errors.Join(ErrInvalidCredentials, errors.New("account with that email does not exist"))
	} else if err != nil {
		return "", errors.Join(ErrFailedToAuth, err)
	}

	if !user.PasswordHash.Valid {
		return "", errors.Join(ErrInvalidCredentials, errors.New("no password set"))
	}
	pwdRaw, err := base64.RawStdEncoding.DecodeString(user.PasswordHash.String)
	if err != nil {
		return "", errors.Join(ErrFailedToAuth, err)
	}
	if err := bcrypt.CompareHashAndPassword(pwdRaw, []byte(password)); err != nil {
		return "", errors.Join(ErrInvalidCredentials, errors.New("passwords do not match"))
	}

	jwt, err = a.createJwt(&user)
	if err != nil {
		return "", errors.Join(ErrFailedToAuth, err)
	}

	return
}

// creates user and returns jwt. Returns [ErrFailedToAuth] if
// something goes wrong [ErrAccountExists] if an account with that email already exists
func (a *Auth) SignUp(email, password string) (jwt string, err error) {
	_, err = a.q.GetUser(context.Background(), email)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return "", errors.Join(ErrFailedToAuth, err)
	} else if err == nil {
		return "", ErrAccountExists
	}

	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.Join(ErrFailedToAuth, err)
	}
	hashedPwd64 := base64.RawStdEncoding.EncodeToString(hashedPwd)

	user, err := a.q.CreateUser(context.Background(), models.CreateUserParams{
		Email:        email,
		PasswordHash: sql.NullString{String: hashedPwd64, Valid: true},
	})
	if err != nil {
		return "", errors.Join(ErrFailedToAuth, err)
	}

	jwt, err = a.createJwt(&user)
	if err != nil {
		return "", errors.Join(ErrFailedToAuth, err)
	}
	return
}

const jwtHeader64 = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"
const expInSec = 86400 * 7 // 7 days

type JwtPayload struct {
	Iat  int64   `json:"iat"`
	Exp  int64   `json:"exp"`
	User JwtUser `json:"user"`
}
type JwtUser struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
}

func (a *Auth) createJwt(user *models.User) (jwt string, err error) {
	nowSec := time.Now().Unix()
	payloadJson, err := json.Marshal(JwtPayload{
		Iat: nowSec,
		Exp: nowSec + expInSec,
		User: JwtUser{
			ID:    user.ID,
			Email: user.Email,
		},
	})
	if err != nil {
		return "", errors.Join(ErrFailedToCreateToken, err)
	}
	payload64 := base64.RawURLEncoding.EncodeToString(payloadJson)

	jwtData := jwtHeader64 + "." + payload64
	fmt.Println(jwtData)

	h := hmac.New(sha256.New, a.secret)
	h.Write([]byte(jwtData))
	signature := base64.RawURLEncoding.EncodeToString(h.Sum(nil))

	return jwtData + "." + signature, nil
}

func (a *Auth) VerifyJwt(jwt string) (user *models.User, err error) {
	s := strings.Split(jwt, ".")
	if len(s) != 3 {
		return nil, errors.Join(ErrInvalidCredentials, fmt.Errorf("expected 3 jwt sections but found %d", len(s)))
	} else if s[0] != jwtHeader64 {
		return nil, errors.Join(ErrInvalidCredentials, errors.New("jwt header does not match"))
	}

	h := hmac.New(sha256.New, a.secret)
	h.Write([]byte(s[0] + "." + s[1]))
	signature := base64.RawURLEncoding.EncodeToString(h.Sum(nil))

	if signature != s[2] {
		return nil, errors.Join(ErrInvalidCredentials, errors.New("signature does not match"))
	}

	payloadJson, err := base64.RawURLEncoding.DecodeString(s[1])
	if err != nil {
		return nil, errors.Join(ErrInvalidCredentials, errors.New("payload is not properly base64 encoded"))
	}
	payload := JwtPayload{}
	err = json.Unmarshal(payloadJson, &payload)
	if err != nil {
		return nil, errors.Join(ErrInvalidCredentials, errors.New("payload is not a valid json format"))
	}

	if time.Now().Unix() > payload.Exp {
		return nil, ErrExpiredCredentials
	}

	dbUser, err := a.q.GetUser(context.Background(), payload.User.Email)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.Join(ErrInvalidCredentials, errors.New("user does not exist"))
	} else if err != nil {
		return nil, errors.Join(err, fmt.Errorf("error fetching DB user %s", payload.User.Email))
	} else if dbUser.ID != payload.User.ID {
		return nil, errors.Join(ErrInvalidCredentials, errors.New("users ids do not match"))
	}
	return &dbUser, nil
}
