package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"slices"
	"strings"
	"time"
	"ultigamecast/internal/app/service/slug"
	"ultigamecast/internal/assert"
	"ultigamecast/internal/ctxvar"
	"ultigamecast/internal/models"
)

type Game struct {
	q  *models.Queries
	db *sql.DB
	clearer SubscriptionClearer
}

type SubscriptionClearer interface{
	ClearGameSubscriptions(ctx context.Context) 
}

func NewGame(db *sql.DB, clearer SubscriptionClearer) *Game {
	return &Game{
		q:  models.New(db),
		db: db,
		clearer: clearer,
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
	tournament := ctxvar.GetTournament(ctx)
	assert.That(tournament != nil, "tournament is nil")
	game, err := g.q.GetGameBySlug(ctx, models.GetGameBySlugParams{
		TournamentID: tournament.ID,
		Slug: strings.ToLower(slug),
	})
	if err != nil {
		return nil, convertAndLogSqlError(ctx, "error getting game", err)
	}
	return &game, nil
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
	if (startTime.Before(tournament.StartDate.Time) || startTime.After(tournament.EndDate.Time.AddDate(0, 0, 1))) && tournament.StartDate.Valid {
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

func (g *Game) UpdateScheduleStatus(ctx context.Context, status string) (*models.Game, error) {
	game := ctxvar.GetGame(ctx)
	assert.That(game != nil, "game cannot be nil when updating schedule status")
	scheduleStatus, ok := models.GameScheduleStatusMap[status]
	assert.That(ok, "unexpected status provided %s", status)

	if game.ScheduleStatus == scheduleStatus {
		return game, nil
	}
	// clear subs if game was live and is now not live
	if scheduleStatus != models.GameScheduleStatusLive && game.ScheduleStatus == models.GameScheduleStatusLive {
		g.clearer.ClearGameSubscriptions(ctx)
	}
	if updatedGame, err := g.q.UpdateGameScheduleStatus(ctx, models.UpdateGameScheduleStatusParams{ID: game.ID, ScheduleStatus: scheduleStatus}); err != nil {
		return nil, convertAndLogSqlError(ctx, "error updating game schedule status", err)
	} else {
		return &updatedGame, nil
	}
}
