package service

import (
	"context"
	"testing"
	"time"
	"ultigamecast/internal/models"
	"ultigamecast/test/testctx"
	"ultigamecast/test/testdb"

	"github.com/stretchr/testify/assert"
)

type clearer struct {
	Count int
}

func (c *clearer) ClearGameSubscriptions(ctx context.Context) {
	c.Count++
}

func TestGame(t *testing.T) {
	locStr := "America/New_York"
	loc, err := time.LoadLocation(locStr)
	if err != nil {
		panic(err)
	}
	startTimeStr := "2024-01-01T15:00"
	validStartTime, err := time.ParseInLocation("2006-01-02T15:04", startTimeStr, loc)
	if err != nil {
		panic(err)
	}
	assert.Equal(t, validStartTime.Format("2006-01-02T15:04"), "2024-01-01T15:00")

	clearerObj := &clearer{}
	game := NewGame(testdb.DB(), clearerObj)
	tournament := NewTournament(testdb.DB())
	teamCtx := testctx.LoadTeam(game.q)
	t.Run("can create game", func(t *testing.T) {
		to, err := tournament.CreateTournament(teamCtx, "tournament", "Jan 1, 2024 - Jan 2, 2024")
		assert.Nil(t, err, "error creating tournament %s", err)
		toCtx := testctx.Load(teamCtx, to)
		game, err := game.CreateGame(toCtx, "opponent", startTimeStr, locStr, 1, 2, 3)
		assert.Nil(t, err, "error creating game %s", err)

		assert.Equal(t, models.GameLiveStatusPrePoint, game.LiveStatus)
		assert.Equal(t, models.GameScheduleStatusScheduled, game.ScheduleStatus)
		gameLoc, err := time.LoadLocation(game.StartTimezone.String)
		assert.Nil(t, err, "error loading game location %s: %s", game.StartTimezone.String, err)
		assert.Equal(t, validStartTime, game.Start.Time.In(gameLoc))
		assert.Equal(t, int64(1), game.HalfCap.Int64)
		assert.Equal(t, int64(2), game.SoftCap.Int64)
		assert.Equal(t, int64(3), game.HardCap.Int64)
	})
	t.Run("cannot create game outside bounds of tournament", func(t *testing.T) {
		to, err := tournament.CreateTournament(teamCtx, "tournament", "Jan 1, 2024 - Jan 2, 2024")
		assert.Nil(t, err, "error creating tournament %s", err)
		toCtx := testctx.Load(teamCtx, to)
		_, err = game.CreateGame(toCtx, "opponent", "2023-01-01T15:00", locStr, 1, 2, 3)
		assert.ErrorIs(t, err, ErrDateOutOfBounds)
	})
	t.Run("can create game with any time when tournament has no date specified", func(t *testing.T) {
		to, err := tournament.CreateTournament(teamCtx, "tournament", "")
		assert.Nil(t, err, "error creating tournament %s", err)
		toCtx := testctx.Load(teamCtx, to)
		_, err = game.CreateGame(toCtx, "opponent", "2023-01-01T15:00", locStr, 1, 2, 3)
		assert.Nil(t, err, "error creating game: %s", err)
	})
	t.Run("removing game from live clears game subscriptions", func(t *testing.T) {
		to, err := tournament.CreateTournament(teamCtx, "tournament", "Jan 1, 2024 - Jan 2, 2024")
		assert.Nil(t, err, "error creating tournament %s", err)
		toCtx := testctx.Load(teamCtx, to)
		g, err := game.CreateGame(toCtx, "opponent", "2024-01-01T15:00", locStr, 1, 2, 3)
		assert.Nil(t, err, "error creating game %s", err)

		clearerObj.Count = 0
		assert.Equal(t, models.GameScheduleStatusScheduled, g.ScheduleStatus)
		gameCtx := testctx.Load(toCtx, g)
		g, err = game.UpdateScheduleStatus(gameCtx, string(models.GameScheduleStatusLive))
		assert.Nil(t, err, "error updating game status %s", err)
		assert.Equal(t, models.GameScheduleStatusLive, g.ScheduleStatus)
		assert.Equal(t, 0, clearerObj.Count)
		gameCtx = testctx.Load(toCtx, g)
		g, err = game.UpdateScheduleStatus(gameCtx, string(models.GameScheduleStatusLive))
		assert.Nil(t, err, "error updating game status %s", err)
		assert.Equal(t, models.GameScheduleStatusLive, g.ScheduleStatus)
		assert.Equal(t, 0, clearerObj.Count)
		gameCtx = testctx.Load(toCtx, g)
		g, err = game.UpdateScheduleStatus(gameCtx, string(models.GameScheduleStatusFinal))
		assert.Nil(t, err, "error updating game status %s", err)
		assert.Equal(t, models.GameScheduleStatusFinal, g.ScheduleStatus)
		assert.Equal(t, 1, clearerObj.Count)
	})
}
