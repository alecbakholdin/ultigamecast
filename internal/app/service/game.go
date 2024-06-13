package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"slices"
	"time"
	"ultigamecast/internal/app/service/slug"
	"ultigamecast/internal/assert"
	"ultigamecast/internal/ctxvar"
	"ultigamecast/internal/models"
)

type Game struct {
	q  *models.Queries
	db *sql.DB
}

func NewGame(q *models.Queries, db *sql.DB) *Game {
	return &Game{
		q:  q,
		db: db,
	}
}

func (g *Game) GetSchedule(ctx context.Context) ([]models.Game, error) {
	tournament := ctxvar.GetTournament(ctx)
	assert.That(tournament != nil, "tournament is not set")
	games, err := g.q.ListTournamentGames(ctx, tournament.ID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, convertAndLogSqlError(ctx, "error fetching tournament games", err)
	}
	return games, nil
}

func (g *Game) GetGame(ctx context.Context, slug string) (*models.Game, error) {
	return nil, nil
}

func (g *Game) CreateGame(ctx context.Context, opponent, start, startTimezone string, half, soft, hard int) (*models.Game, error){
	tournament := ctxvar.GetTournament(ctx)
	assert.That(tournament != nil, "tournament is not set")
	slug, err := g.getSafeSlug(ctx, opponent)
	if err != nil {
		return nil, convertAndLogSqlError(ctx, "error generating slug", err)
	}

	loc, err := time.LoadLocation(startTimezone)
	if err != nil {
		slog.ErrorContext(ctx, "error loading location", "location", startTimezone, "err", err)
		return nil, err
	}
	var startTime time.Time
	if start != "" {
		startTime, err = time.ParseInLocation("2006-01-02T15:04", start, loc)
		if err != nil {
			slog.ErrorContext(ctx, "error parsing time", "timestr", start, "err", err)
			return nil, err
		}
	}
	if startTime.Before(tournament.StartDate.Time) || startTime.After(tournament.EndDate.Time.AddDate(0, 0, 1)) {
		return nil, ErrDateOutOfBounds
	}

	game, err := g.q.CreateGame(ctx, models.CreateGameParams{
		Tournament: tournament.ID,
		Opponent: opponent,
		Slug: slug,

		Start: sql.NullTime{Time: startTime, Valid: !startTime.IsZero()},
		StartTimezone: sql.NullString{String: startTimezone, Valid: startTimezone != "" },

		HalfCap: sql.NullInt64{Int64: int64(half), Valid: half > 0},
		SoftCap: sql.NullInt64{Int64: int64(soft), Valid: soft > 0},
		HardCap: sql.NullInt64{Int64: int64(hard), Valid: hard > 0},
	})
	if err != nil {
		return nil, convertAndLogSqlError(ctx, "error creating game", err)
	}
	return &game, nil
}

func (t *Game) getSafeSlug(ctx context.Context, opponent string) (string, error) {
	tournament := ctxvar.GetTournament(ctx)
	tournaments, err := t.q.ListTournamentGames(ctx, tournament.ID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return "", err
	}
	s := slug.From(opponent)
	num := 2
	for slices.ContainsFunc(tournaments, func(to models.Game) bool { return to.ID != tournament.ID && to.Slug == s }) {
		s = fmt.Sprintf("%s-%d", slug.From(opponent), num)
		num++
	}
	return s, nil
}