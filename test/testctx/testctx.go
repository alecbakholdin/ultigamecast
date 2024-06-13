package testctx

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"
	"strings"
	"ultigamecast/internal/ctxvar"
	"ultigamecast/internal/models"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

const randChars = "abcdefghjiklmnopqrstuvwxyz0123456789"

func randomString(n int) string {
	str := ""
	for range n {
		str += string(randChars[rand.Int()%len(randChars)])
	}
	return str
}

func LoadUser(q *models.Queries) context.Context {
	ctx := context.Background()
	pwd, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		panic(fmt.Errorf("error bcrypting pwd: %w", err))
	}
	user, err := q.CreateUser(ctx, models.CreateUserParams{
		Email:        fmt.Sprintf("user%s@gmail.com", randomString(10)),
		PasswordHash: sql.NullString{String: string(pwd), Valid: true},
	})
	if err != nil {
		panic(fmt.Errorf("error creating user: %w", err))
	}
	return Load(ctx, user)
}

func LoadTeam(q *models.Queries) context.Context {
	ctx := LoadUser(q)
	randstr := randomString(10)
	team, err := q.CreateTeam(ctx, models.CreateTeamParams{
		Owner: ctxvar.GetUser(ctx).ID,
		Name:  "team " + randstr,
		Slug:  "team-" + strings.ToLower(randstr),
	})
	if err != nil {
		panic(fmt.Errorf("error creating team: %w", err))
	}
	return Load(ctx, team)
}

func LoadTournament(q *models.Queries) context.Context {
	ctx := LoadTeam(q)
	randstr := randomString(10)
	tournament, err := q.CreateTournament(ctx, models.CreateTournamentParams{
		TeamId: ctxvar.GetTeam(ctx).ID,
		Name: "tournament " + randstr,
		Slug: "tournament-" + strings.ToLower(randstr),
	})
	if err != nil {
		panic(fmt.Errorf("error creating tournament: %w", err))
	}
	return Load(ctx, tournament)
}

func LoadTournamentDatum(q *models.Queries) context.Context {
	ctx := LoadTournament(q)
	datum, err := q.CreateTournamentDatum(ctx, ctxvar.GetTournament(ctx).ID)
	if err != nil {
		panic(fmt.Errorf("error creating datum: %w", err))
	}
	return Load(ctx, datum)
}

func Load(ctx context.Context, values ...any) context.Context {
	for _, val := range values {
		switch v := val.(type) {
		case models.User:
			ctx = context.WithValue(ctx, ctxvar.User, &v)
		case *models.User:
			ctx = context.WithValue(ctx, ctxvar.User, v)
		case models.Team:
			ctx = context.WithValue(ctx, ctxvar.Team, &v)
		case *models.Team:
			ctx = context.WithValue(ctx, ctxvar.Team, v)
		case models.Tournament:
			ctx = context.WithValue(ctx, ctxvar.Tournament, &v)
		case *models.Tournament:
			ctx = context.WithValue(ctx, ctxvar.Tournament, v)
		case models.TournamentSummary:
			ctx = context.WithValue(ctx, ctxvar.Tournament, v.Tournament)
		case *models.TournamentSummary:
			ctx = context.WithValue(ctx, ctxvar.Tournament, v.Tournament)
		case models.TournamentDatum:
			ctx = context.WithValue(ctx, ctxvar.TournamentDatum, &v)
		case *models.TournamentDatum:
			ctx = context.WithValue(ctx, ctxvar.TournamentDatum, v)
		case models.Player:
			ctx = context.WithValue(ctx, ctxvar.Player, &v)
		case *models.Player:
			ctx = context.WithValue(ctx, ctxvar.Player, v)
		case models.Game:
			ctx = context.WithValue(ctx, ctxvar.Game, &v)
		case *models.Game:
			ctx = context.WithValue(ctx, ctxvar.Game, v)
		case models.Event:
			ctx = context.WithValue(ctx, ctxvar.Event, &v)
		case *models.Event:
			ctx = context.WithValue(ctx, ctxvar.Event, v)
		default:
			panic(fmt.Errorf("unexpected type while loading context %T", v))
		}
	}
	return ctx
}
