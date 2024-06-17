package service

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
	"ultigamecast/internal/models"
	"ultigamecast/test/testctx"
	"ultigamecast/test/testdb"

	"github.com/stretchr/testify/assert"
)

func TestEvent(t *testing.T) {
	event := NewEvent(testdb.DB())
	tournament := NewTournament(testdb.DB())
	game := NewGame(testdb.DB(), event)
	player := NewPlayer(testdb.DB())
	teamCtx := testctx.LoadTeam(event.q)
	playerSlugs := []string{}
	playerIds := []int64{}
	playerIdStrings := []string{}
	for i := range 7 {
		if p, err := player.CreatePlayer(teamCtx, fmt.Sprintf("player %d", i)); err != nil {
			t.Fatalf("error creating player %d: %s", i, err)
		} else {
			playerSlugs = append(playerSlugs, p.Slug)
			playerIds = append(playerIds, p.ID)
			playerIdStrings = append(playerIdStrings, strconv.FormatInt(p.ID, 10))
		}
	}
	t.Run("subscription fails if game is not live", func(t *testing.T) {
		// create tournament and game
		to, err := tournament.CreateTournament(teamCtx, "tournament", "Jan 2, 2024 - Jan 3, 2024")
		assert.Nil(t, err, "error creating tournament")
		toCtx := testctx.Load(teamCtx, to)
		g, err := game.CreateGame(toCtx, "opponent", "2024-01-02T15:00", "America/New_York", 0, 0, 0)
		assert.Nil(t, err, "error creating game: %s", err)

		// test
		_, err = event.Subscribe(testctx.Load(toCtx, g))
		assert.ErrorIs(t, err, ErrGameNotLive)
	})
	t.Run("general game flow works", func(t *testing.T) {
		// create tournament and game
		to, err := tournament.CreateTournament(teamCtx, "tournament", "Jan 2, 2024 - Jan 3, 2024")
		assert.Nil(t, err, "error creating tournament")
		toCtx := testctx.Load(teamCtx, to)
		g, err := game.CreateGame(toCtx, "opponent", "2024-01-02T15:00", "America/New_York", 0, 0, 0)
		assert.Nil(t, err, "error creating game: %s", err)
		gameCtx := testctx.Load(toCtx, g)
		g, err = game.UpdateScheduleStatus(gameCtx, string(models.GameScheduleStatusLive))
		assert.Nil(t, err, "error updating game status: %s", err)
		gameCtx = testctx.Load(toCtx, g)

		// subscribe
		sub, err := event.Subscribe(gameCtx)
		assert.Nil(t, err, "error subscribing: %s", err)
		assert.Equal(t, g.ID, sub.GameId)
		assert.NotNil(t, sub.Id)
		assert.NotNil(t, sub.EventChan)

		t.Run("01 starting line", func(t *testing.T) {
			err = event.StartingLine(gameCtx, playerSlugs...)
			assert.Nil(t, err, "%s", err)
			if err != nil {
				return
			}
			var batch string
			for range 7 {
				u := <-sub.EventChan

				if batch == "" {
					batch = u.Event.Batch.String
				}
				assert.NotEmpty(t, batch)
				assert.Equal(t, batch, u.Event.Batch.String)
				assert.Equal(t, batch, u.Game.LastEvent.String)

				assert.Equal(t, models.EventTypeStartingLine, u.Event.Type)
				assert.Equal(t, models.GameLiveStatusPrePoint, u.Game.LiveStatus)
				assert.True(t, u.Game.ActivePlayers.Valid)
				assert.Equal(t, strings.Join(playerIdStrings, ","), u.Game.ActivePlayers.String)
			}
		})

		t.Run("02 start point defense", func(t *testing.T) {
			err = event.StartPointDefense(gameCtx)
			assert.Nil(t, err)
			if err != nil {
				return
			}
			u := <-sub.EventChan
			assert.Equal(t, models.EventTypePointStart, u.Event.Type)
			assert.Equal(t, models.GameLiveStatusOpponentPossession, u.Game.LiveStatus)
			assert.Equal(t, int64(0), u.Event.TeamScore)
			assert.Equal(t, int64(0), u.Event.OpponentScore)
			assert.Equal(t, u.Event.ID, u.Game.LastEvent.String)
		})

		t.Run("03 opponent drop", func(t *testing.T) {
			err = event.OpponentDrop(gameCtx)
			assert.Nil(t, err)
			if err != nil {
				return
			}
			u := <-sub.EventChan
			assert.Equal(t, models.EventTypeOpponentDrop, u.Event.Type)
			assert.Equal(t, models.GameLiveStatusTeamPossession, u.Game.LiveStatus)
			assert.Equal(t, int64(0), u.Event.TeamScore)
			assert.Equal(t, int64(0), u.Event.OpponentScore)
			assert.Equal(t, u.Event.ID, u.Game.LastEvent.String)
		})

		t.Run("04 team drop", func(t *testing.T) {
			err = event.Drop(gameCtx, playerSlugs[0])
			assert.Nil(t, err)
			if err != nil {
				return
			}
			u := <-sub.EventChan
			assert.Equal(t, playerIds[0], u.Event.Player.Int64)
			assert.Equal(t, models.EventTypeDrop, u.Event.Type)
			assert.Equal(t, models.GameLiveStatusOpponentPossession, u.Game.LiveStatus)
			assert.Equal(t, int64(0), u.Event.TeamScore)
			assert.Equal(t, int64(0), u.Event.OpponentScore)
			assert.Equal(t, u.Event.ID, u.Game.LastEvent.String)
		})

		t.Run("05 opponent timeout", func(t *testing.T) {
			err = event.OpponentTimeoutPoint(gameCtx)
			assert.Nil(t, err)
			if err != nil {
				return
			}
			u := <-sub.EventChan
			assert.Equal(t, models.EventTypeOpponentTimeout, u.Event.Type)
			assert.Equal(t, models.GameLiveStatusOpponentTimeout, u.Game.LiveStatus)
			assert.Equal(t, int64(0), u.Event.TeamScore)
			assert.Equal(t, int64(0), u.Event.OpponentScore)
			assert.Equal(t, u.Event.ID, u.Game.LastEvent.String)
		})

		t.Run("06 opponent end timeout", func(t *testing.T) {
			err = event.OpponentTimeoutEndPoint(gameCtx)
			assert.Nil(t, err)
			if err != nil {
				return
			}
			u := <-sub.EventChan
			assert.Equal(t, models.EventTypeOpponentTimeoutEnd, u.Event.Type)
			assert.Equal(t, models.GameLiveStatusOpponentPossession, u.Game.LiveStatus)
			assert.Equal(t, int64(0), u.Event.TeamScore)
			assert.Equal(t, int64(0), u.Event.OpponentScore)
			assert.Equal(t, u.Event.ID, u.Game.LastEvent.String)
		})

		t.Run("07 opponent turn", func(t *testing.T) {
			err = event.OpponentTurn(gameCtx)
			assert.Nil(t, err)
			if err != nil {
				return
			}
			u := <-sub.EventChan
			assert.Equal(t, models.EventTypeOpponentTurn, u.Event.Type)
			assert.Equal(t, models.GameLiveStatusTeamPossession, u.Game.LiveStatus)
			assert.Equal(t, int64(0), u.Event.TeamScore)
			assert.Equal(t, int64(0), u.Event.OpponentScore)
			assert.Equal(t, u.Event.ID, u.Game.LastEvent.String)
		})

		t.Run("08 team timeout", func(t *testing.T) {
			err = event.TimeoutPoint(gameCtx)
			assert.Nil(t, err)
			if err != nil {
				return
			}
			u := <-sub.EventChan
			assert.Equal(t, models.EventTypeTimeout, u.Event.Type)
			assert.Equal(t, models.GameLiveStatusTeamTimeout, u.Game.LiveStatus)
			assert.Equal(t, int64(0), u.Event.TeamScore)
			assert.Equal(t, int64(0), u.Event.OpponentScore)
			assert.Equal(t, u.Event.ID, u.Game.LastEvent.String)
		})

		t.Run("09 team end timeout", func(t *testing.T) {
			err = event.TimeoutEndPoint(gameCtx)
			assert.Nil(t, err)
			if err != nil {
				return
			}
			u := <-sub.EventChan
			assert.Equal(t, models.EventTypeTimeoutEnd, u.Event.Type)
			assert.Equal(t, models.GameLiveStatusTeamPossession, u.Game.LiveStatus)
			assert.Equal(t, int64(0), u.Event.TeamScore)
			assert.Equal(t, int64(0), u.Event.OpponentScore)
			assert.Equal(t, u.Event.ID, u.Game.LastEvent.String)
		})

		t.Run("10 team goal", func(t *testing.T) {
			err = event.Goal(gameCtx, playerSlugs[0], playerSlugs[1])
			assert.Nil(t, err)
			if err != nil {
				return
			}
			u := <-sub.EventChan
			assert.Equal(t, models.EventTypeGoal, u.Event.Type)
			assert.Equal(t, playerIds[0], u.Event.Player.Int64)
			assert.Equal(t, models.GameLiveStatusPrePoint, u.Game.LiveStatus)
			assert.Equal(t, int64(0), u.Event.TeamScore)
			assert.Equal(t, int64(0), u.Event.OpponentScore)
			assert.Equal(t, u.Event.Batch.String, u.Game.LastEvent.String)

			u = <-sub.EventChan
			assert.Equal(t, models.EventTypeAssist, u.Event.Type)
			assert.Equal(t, playerIds[1], u.Event.Player.Int64)
			assert.Equal(t, models.GameLiveStatusPrePoint, u.Game.LiveStatus)
			assert.Equal(t, int64(0), u.Event.TeamScore)
			assert.Equal(t, int64(0), u.Event.OpponentScore)
			assert.Equal(t, u.Event.Batch.String, u.Game.LastEvent.String)

			assert.Equal(t, int64(1), u.Game.TeamScore)
			assert.Equal(t, int64(0), u.Game.OpponentScore)
		})

		t.Run("11 opponent prepoint timeout", func(t *testing.T) {
			err = event.OpponentTimeout(gameCtx)
			assert.Nil(t, err)
			if err != nil {
				return
			}
			u := <-sub.EventChan
			assert.Equal(t, models.EventTypeOpponentTimeout, u.Event.Type)
			assert.Equal(t, models.GameLiveStatusOpponentTimeout, u.Game.LiveStatus)
			assert.Equal(t, int64(1), u.Event.TeamScore)
			assert.Equal(t, int64(0), u.Event.OpponentScore)
			assert.Equal(t, u.Event.ID, u.Game.LastEvent.String)
		})

		t.Run("12 opponent prepoint timeout end", func(t *testing.T) {
			err = event.OpponentTimeoutEnd(gameCtx)
			assert.Nil(t, err)
			if err != nil {
				return
			}
			u := <-sub.EventChan
			assert.Equal(t, models.EventTypeOpponentTimeoutEnd, u.Event.Type)
			assert.Equal(t, models.GameLiveStatusPrePoint, u.Game.LiveStatus)
			assert.Equal(t, int64(1), u.Event.TeamScore)
			assert.Equal(t, int64(0), u.Event.OpponentScore)
			assert.Equal(t, u.Event.ID, u.Game.LastEvent.String)
		})

		t.Run("13 team prepoint timeout", func(t *testing.T) {
			err = event.Timeout(gameCtx)
			assert.Nil(t, err)
			if err != nil {
				return
			}
			u := <-sub.EventChan
			assert.Equal(t, models.EventTypeTimeout, u.Event.Type)
			assert.Equal(t, models.GameLiveStatusPrePoint, u.Game.LiveStatus)
			assert.Equal(t, int64(1), u.Event.TeamScore)
			assert.Equal(t, int64(0), u.Event.OpponentScore)
			assert.Equal(t, u.Event.ID, u.Game.LastEvent.String)
		})

		t.Run("14 team prepoint timeout end", func(t *testing.T) {
			err = event.TimeoutEnd(gameCtx)
			assert.Nil(t, err)
			if err != nil {
				return
			}
			u := <-sub.EventChan
			assert.Equal(t, models.EventTypeTimeoutEnd, u.Event.Type)
			assert.Equal(t, models.GameLiveStatusPrePoint, u.Game.LiveStatus)
			assert.Equal(t, int64(1), u.Event.TeamScore)
			assert.Equal(t, int64(0), u.Event.OpponentScore)
			assert.Equal(t, u.Event.ID, u.Game.LastEvent.String)
		})

		t.Run("15 start point defense", func(t *testing.T) {
			err = event.StartingLine(gameCtx, playerSlugs...)
			assert.Nil(t, err)
			for range 7 {
				<-sub.EventChan
			}
			if err != nil {
				return
			}
			err = event.StartPointDefense(gameCtx)
			assert.Nil(t, err)
			if err != nil {
				return
			}
			u := <-sub.EventChan
			assert.Equal(t, models.EventTypePointStart, u.Event.Type)
			assert.Equal(t, models.GameLiveStatusOpponentPossession, u.Game.LiveStatus)
			assert.Equal(t, int64(1), u.Event.TeamScore)
			assert.Equal(t, int64(0), u.Event.OpponentScore)
			assert.Equal(t, u.Event.ID, u.Game.LastEvent.String)
		})

		t.Run("16 opponent goal", func(t *testing.T) {
			err = event.OpponentGoal(gameCtx)
			assert.Nil(t, err)
			if err != nil {
				return
			}
			u := <-sub.EventChan
			assert.Equal(t, models.EventTypeOpponentGoal, u.Event.Type)
			assert.Equal(t, models.GameLiveStatusPrePoint, u.Game.LiveStatus)
			assert.Equal(t, int64(1), u.Event.TeamScore)
			assert.Equal(t, int64(0), u.Event.OpponentScore)
			assert.Equal(t, u.Event.ID, u.Game.LastEvent.String)

			assert.Equal(t, int64(1), u.Game.TeamScore)
			assert.Equal(t, int64(1), u.Game.OpponentScore)
		})

		t.Run("17 start point offense", func(t *testing.T) {
			err = event.StartingLine(gameCtx, playerSlugs...)
			for range 7 {
				<-sub.EventChan
			}
			assert.Nil(t, err)
			if err != nil {
				return
			}
			err = event.StartPointOffense(gameCtx)
			assert.Nil(t, err)
			if err != nil {
				return
			}
			u := <-sub.EventChan
			assert.Equal(t, models.EventTypePointStart, u.Event.Type)
			assert.Equal(t, models.GameLiveStatusTeamPossession, u.Game.LiveStatus)
			assert.Equal(t, int64(1), u.Event.TeamScore)
			assert.Equal(t, int64(1), u.Event.OpponentScore)
			assert.Equal(t, u.Event.ID, u.Game.LastEvent.String)
		})

		t.Run("18 team goal", func(t *testing.T) {
			err = event.Goal(gameCtx, playerSlugs[1], playerSlugs[0])
			assert.Nil(t, err)
			if err != nil {
				return
			}
			u := <-sub.EventChan
			assert.Equal(t, playerIds[1], u.Event.Player.Int64)
			assert.Equal(t, models.EventTypeGoal, u.Event.Type)
			assert.Equal(t, models.GameLiveStatusPrePoint, u.Game.LiveStatus)
			assert.Equal(t, int64(1), u.Event.TeamScore)
			assert.Equal(t, int64(1), u.Event.OpponentScore)

			u = <-sub.EventChan
			assert.Equal(t, playerIds[0], u.Event.Player.Int64)
			assert.Equal(t, models.EventTypeAssist, u.Event.Type)
			assert.Equal(t, models.GameLiveStatusPrePoint, u.Game.LiveStatus)
			assert.Equal(t, int64(1), u.Event.TeamScore)
			assert.Equal(t, int64(1), u.Event.OpponentScore)

			assert.Equal(t, u.Event.Batch, u.Game.LastEvent)
			assert.Equal(t, int64(2), u.Game.TeamScore)
			assert.Equal(t, int64(1), u.Game.OpponentScore)
		})
	})
}
