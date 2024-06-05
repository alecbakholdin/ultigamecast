package testctx

import (
	"context"
	"fmt"
	"ultigamecast/internal/ctxvar"
	"ultigamecast/internal/models"

	_ "github.com/mattn/go-sqlite3"
)

func LoadUser(q *models.Queries) context.Context {
	if user, err := q.GetUser(context.Background(), "alecbakholdin@gmail.com"); err != nil {
		panic(fmt.Errorf("error fetching user: %w", err))
	} else {
		return Load(context.Background(), &user)
	}
}

func LoadTeam(q *models.Queries) context.Context {
	ctx := LoadUser(q)
	if teams, err := q.ListOwnedTeams(ctx, ctxvar.GetUser(ctx).ID); err != nil {
		panic(fmt.Errorf("error fetching teams: %w", err))
	} else {
		return Load(ctx, &teams[0])
	}
}

func LoadTournament(q *models.Queries) context.Context {
	ctx := LoadTeam(q)
	if tournaments, err := q.ListTournaments(ctx, ctxvar.GetTeam(ctx).ID); err != nil {
		panic(fmt.Errorf("error fetching teams: %w", err))
	} else {
		return Load(ctx, &tournaments[0])
	}
}

func LoadTournamentDatum(q *models.Queries) context.Context {
	ctx := LoadTournament(q)
	if data, err := q.ListTournamentData(ctx, ctxvar.GetTournament(ctx).ID); err != nil {
		panic(fmt.Errorf("error fetching data: %w", err))
	} else {
		return Load(ctx, &data[0])
	}
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
			panic(fmt.Errorf("unexpected type %T", v))
		}
	}
	return ctx
}