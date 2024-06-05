package testdb

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"ultigamecast/internal/ctxvar"
	"ultigamecast/internal/models"

	_ "github.com/mattn/go-sqlite3"
)

var (
	_, b, _, _ = runtime.Caller(0)
	basepath   = filepath.Dir(b)
	testDbPath = filepath.Join(basepath, "..", "test.db")
)
var testdb *sql.DB

func DB() (*models.Queries, *sql.DB) {
	if testdb == nil {
		Init()
	}
	return models.New(testdb), testdb
}

func Init() {
	ogTestDb, err := os.Open(testDbPath)
	if err != nil {
		panic(fmt.Errorf("error opening file: %w", err))
	}
	defer ogTestDb.Close()
	newTestDb, err := os.CreateTemp(os.TempDir(), "*.db")
	if err != nil {
		panic(fmt.Errorf("error creating temp file: %w", err))
	}

	_, err = io.Copy(newTestDb, ogTestDb)
	newTestDb.Close()
	if err != nil {
		panic(fmt.Errorf("error copying to new file: %w", err))
	}

	testdb, err = sql.Open("sqlite3", newTestDb.Name())
	if err != nil {
		panic(err)
	}
}

func LoadUser(q *models.Queries) context.Context {
	if user, err := q.GetUser(context.Background(), "alecbakholdin@gmail.com"); err != nil {
		panic(fmt.Errorf("error fetching user: %w", err))
	} else {
		return LoadCtxValue(context.Background(), &user)
	}
}

func LoadTeam(q *models.Queries) context.Context {
	ctx := LoadUser(q)
	if teams, err := q.ListOwnedTeams(ctx, ctxvar.GetUser(ctx).ID); err != nil {
		panic(fmt.Errorf("error fetching teams: %w", err))
	} else {
		return LoadCtxValue(ctx, &teams[0])
	}
}

func LoadTournament(q *models.Queries) context.Context {
	ctx := LoadTeam(q)
	if tournaments, err := q.ListTournaments(ctx, ctxvar.GetTeam(ctx).ID); err != nil {
		panic(fmt.Errorf("error fetching teams: %w", err))
	} else {
		return LoadCtxValue(ctx, &tournaments[0])
	}
}

func LoadTournamentDatum(q *models.Queries) context.Context {
	ctx := LoadTournament(q)
	if data, err := q.ListTournamentData(ctx, ctxvar.GetTournament(ctx).ID); err != nil {
		panic(fmt.Errorf("error fetching data: %w", err))
	} else {
		return LoadCtxValue(ctx, &data[0])
	}
}

func LoadCtxValue(ctx context.Context, values ...any) context.Context {
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
