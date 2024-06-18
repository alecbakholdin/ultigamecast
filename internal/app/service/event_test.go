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
			u := <-sub.EventChan
			assert.Equal(t, u.Id, u.Game.LastEvent.String)
			assert.Equal(t, models.GameLiveStatusPrePoint, u.Game.LiveStatus)
			assert.True(t, u.Game.ActivePlayers.Valid)
			assert.Equal(t, strings.Join(playerIdStrings, ","), u.Game.ActivePlayers.String)
			assert.NotEmpty(t, u.Id)
			for i, e := range u.Event{
				assert.Equal(t, u.Id, e.Batch.String)
				assert.Equal(t, models.EventTypeStartingLine, e.Type)
				assert.Equal(t, playerIds[i], e.Player.Int64)
			}
		})

		t.Run("02 start point defense", func(t *testing.T) {
			err = event.StartPointDefense(gameCtx)
			assert.Nil(t, err)
			if err != nil {
				return
			}
			u := <-sub.EventChan
			assert.Equal(t, u.Event[0].ID, u.Game.LastEvent.String)
			assert.Equal(t, models.GameLiveStatusOpponentPossession, u.Game.LiveStatus)
			assert.Equal(t, 1, len(u.Event))
			assert.Equal(t, models.EventTypePointStart, u.Event[0].Type)
			assert.Equal(t, int64(0), u.Event[0].TeamScore)
			assert.Equal(t, int64(0), u.Event[0].OpponentScore)
		})

		t.Run("03 opponent drop", func(t *testing.T) {
			err = event.OpponentDrop(gameCtx)
			assert.Nil(t, err)
			if err != nil {
				return
			}
			u := <-sub.EventChan
			assert.Equal(t, models.GameLiveStatusTeamPossession, u.Game.LiveStatus)
			assert.Equal(t, u.Id, u.Game.LastEvent.String)
			assert.Equal(t, models.EventTypeOpponentDrop, u.Event[0].Type)
			assert.Equal(t, int64(0), u.Event[0].TeamScore)
			assert.Equal(t, int64(0), u.Event[0].OpponentScore)
		})

		t.Run("04 team drop", func(t *testing.T) {
			err = event.Drop(gameCtx, playerSlugs[0])
			assert.Nil(t, err)
			if err != nil {
				return
			}
			u := <-sub.EventChan
			assert.Equal(t, 1, len(u.Event))
			assert.Equal(t, playerIds[0], u.Event[0].Player.Int64)
			assert.Equal(t, u.Id, u.Game.LastEvent.String)
			assert.Equal(t, models.EventTypeDrop, u.Event[0].Type)
			assert.Equal(t, models.GameLiveStatusOpponentPossession, u.Game.LiveStatus)
			assert.Equal(t, int64(0), u.Event[0].TeamScore)
			assert.Equal(t, int64(0), u.Event[0].OpponentScore)
		})

		t.Run("05 opponent timeout", func(t *testing.T) {
			err = event.OpponentTimeoutPoint(gameCtx)
			assert.Nil(t, err)
			if err != nil {
				return
			}
			u := <-sub.EventChan
			assert.Equal(t, 1, len(u.Event))
			assert.Equal(t, u.Id, u.Game.LastEvent.String)
			assert.Equal(t, models.GameLiveStatusOpponentTimeout, u.Game.LiveStatus)
			assert.Equal(t, models.EventTypeOpponentTimeout, u.Event[0].Type)
			assert.Equal(t, int64(0), u.Event[0].TeamScore)
			assert.Equal(t, int64(0), u.Event[0].OpponentScore)
		})

		t.Run("06 opponent end timeout", func(t *testing.T) {
			err = event.OpponentTimeoutEndPoint(gameCtx)
			assert.Nil(t, err)
			if err != nil {
				return
			}
			u := <-sub.EventChan
			assert.Equal(t, 1, len(u.Event))
			assert.Equal(t, u.Id, u.Game.LastEvent.String)
			assert.Equal(t, models.GameLiveStatusOpponentPossession, u.Game.LiveStatus)
			assert.Equal(t, models.EventTypeOpponentTimeoutEnd, u.Event[0].Type)
			assert.Equal(t, int64(0), u.Event[0].TeamScore)
			assert.Equal(t, int64(0), u.Event[0].OpponentScore)
		})

		t.Run("07 opponent turn", func(t *testing.T) {
			err = event.OpponentTurn(gameCtx)
			assert.Nil(t, err)
			if err != nil {
				return
			}
			u := <-sub.EventChan
			assert.Equal(t, 1, len(u.Event))
			assert.Equal(t, u.Id, u.Game.LastEvent.String)
			assert.Equal(t, models.GameLiveStatusTeamPossession, u.Game.LiveStatus)
			assert.Equal(t, models.EventTypeOpponentTurn, u.Event[0].Type)
			assert.Equal(t, int64(0), u.Event[0].TeamScore)
			assert.Equal(t, int64(0), u.Event[0].OpponentScore)
		})

		t.Run("08 team timeout", func(t *testing.T) {
			err = event.TimeoutPoint(gameCtx)
			assert.Nil(t, err)
			if err != nil {
				return
			}
			u := <-sub.EventChan
			assert.Equal(t, 1, len(u.Event))
			assert.Equal(t, u.Id, u.Game.LastEvent.String)
			assert.Equal(t, models.GameLiveStatusTeamTimeout, u.Game.LiveStatus)
			assert.Equal(t, models.EventTypeTimeout, u.Event[0].Type)
			assert.Equal(t, int64(0), u.Event[0].TeamScore)
			assert.Equal(t, int64(0), u.Event[0].OpponentScore)
		})

		t.Run("09 team end timeout", func(t *testing.T) {
			err = event.TimeoutEndPoint(gameCtx)
			assert.Nil(t, err)
			if err != nil {
				return
			}
			u := <-sub.EventChan
			assert.Equal(t, 1, len(u.Event))
			assert.Equal(t, u.Id, u.Game.LastEvent.String)
			assert.Equal(t, models.GameLiveStatusTeamPossession, u.Game.LiveStatus)
			assert.Equal(t, models.EventTypeTimeoutEnd, u.Event[0].Type)
			assert.Equal(t, int64(0), u.Event[0].TeamScore)
			assert.Equal(t, int64(0), u.Event[0].OpponentScore)
		})

		t.Run("10 team goal", func(t *testing.T) {
			err = event.Goal(gameCtx, playerSlugs[0], playerSlugs[1])
			assert.Nil(t, err)
			if err != nil {
				return
			}
			u := <-sub.EventChan
			assert.Equal(t, 2, len(u.Event))
			assert.Equal(t, u.Id, u.Game.LastEvent.String)
			assert.Equal(t, models.GameLiveStatusPrePoint, u.Game.LiveStatus)
			assert.Equal(t, int64(1), u.Game.TeamScore)
			assert.Equal(t, int64(0), u.Game.OpponentScore)

			assert.Equal(t, models.EventTypeGoal, u.Event[0].Type)
			assert.Equal(t, playerIds[0], u.Event[0].Player.Int64)
			assert.Equal(t, int64(0), u.Event[0].TeamScore)
			assert.Equal(t, int64(0), u.Event[0].OpponentScore)
			assert.Equal(t, u.Id, u.Event[0].Batch.String)

			assert.Equal(t, models.EventTypeAssist, u.Event[1].Type)
			assert.Equal(t, playerIds[1], u.Event[1].Player.Int64)
			assert.Equal(t, int64(0), u.Event[1].TeamScore)
			assert.Equal(t, int64(0), u.Event[1].OpponentScore)
			assert.Equal(t, u.Id, u.Event[1].Batch.String)
		})

		t.Run("11 opponent prepoint timeout", func(t *testing.T) {
			err = event.OpponentTimeout(gameCtx)
			assert.Nil(t, err)
			if err != nil {
				return
			}
			u := <-sub.EventChan
			assert.Equal(t, 1, len(u.Event))
			assert.Equal(t, u.Id, u.Game.LastEvent.String)
			assert.Equal(t, models.GameLiveStatusOpponentTimeout, u.Game.LiveStatus)
			assert.Equal(t, models.EventTypeOpponentTimeout, u.Event[0].Type)
			assert.Equal(t, int64(1), u.Event[0].TeamScore)
			assert.Equal(t, int64(0), u.Event[0].OpponentScore)
		})

		t.Run("12 opponent prepoint timeout end", func(t *testing.T) {
			err = event.OpponentTimeoutEnd(gameCtx)
			assert.Nil(t, err)
			if err != nil {
				return
			}
			u := <-sub.EventChan
			assert.Equal(t, 1, len(u.Event))
			assert.Equal(t, u.Id, u.Game.LastEvent.String)
			assert.Equal(t, models.GameLiveStatusPrePoint, u.Game.LiveStatus)
			assert.Equal(t, models.EventTypeOpponentTimeoutEnd, u.Event[0].Type)
			assert.Equal(t, int64(1), u.Event[0].TeamScore)
			assert.Equal(t, int64(0), u.Event[0].OpponentScore)
		})

		t.Run("13 team prepoint timeout", func(t *testing.T) {
			err = event.Timeout(gameCtx)
			assert.Nil(t, err)
			if err != nil {
				return
			}
			u := <-sub.EventChan
			assert.Equal(t, 1, len(u.Event))
			assert.Equal(t, u.Id, u.Game.LastEvent.String)
			assert.Equal(t, models.GameLiveStatusPrePoint, u.Game.LiveStatus)
			assert.Equal(t, models.EventTypeTimeout, u.Event[0].Type)
			assert.Equal(t, int64(1), u.Event[0].TeamScore)
			assert.Equal(t, int64(0), u.Event[0].OpponentScore)
		})

		t.Run("14 team prepoint timeout end", func(t *testing.T) {
			err = event.TimeoutEnd(gameCtx)
			assert.Nil(t, err)
			if err != nil {
				return
			}
			u := <-sub.EventChan
			assert.Equal(t, 1, len(u.Event))
			assert.Equal(t, u.Id, u.Game.LastEvent.String)
			assert.Equal(t, models.GameLiveStatusPrePoint, u.Game.LiveStatus)
			assert.Equal(t, models.EventTypeTimeoutEnd, u.Event[0].Type)
			assert.Equal(t, int64(1), u.Event[0].TeamScore)
			assert.Equal(t, int64(0), u.Event[0].OpponentScore)
		})

		t.Run("15 start point defense", func(t *testing.T) {
			err = event.StartingLine(gameCtx, playerSlugs...)
			assert.Nil(t, err)
			if err != nil {
				return
			}
			<- sub.EventChan
			err = event.StartPointDefense(gameCtx)
			assert.Nil(t, err)
			if err != nil {
				return
			}
			u := <-sub.EventChan
			assert.Equal(t, 1, len(u.Event))
			assert.Equal(t, u.Id, u.Game.LastEvent.String)
			assert.Equal(t, models.GameLiveStatusOpponentPossession, u.Game.LiveStatus)
			assert.Equal(t, models.EventTypePointStart, u.Event[0].Type)
			assert.Equal(t, int64(1), u.Event[0].TeamScore)
			assert.Equal(t, int64(0), u.Event[0].OpponentScore)
		})

		t.Run("16 opponent goal", func(t *testing.T) {
			err = event.OpponentGoal(gameCtx)
			assert.Nil(t, err)
			if err != nil {
				return
			}
			u := <-sub.EventChan
			assert.Equal(t, 1, len(u.Event))
			assert.Equal(t, u.Id, u.Game.LastEvent.String)
			assert.Equal(t, models.GameLiveStatusPrePoint, u.Game.LiveStatus)
			assert.Equal(t, int64(1), u.Game.TeamScore)
			assert.Equal(t, int64(1), u.Game.OpponentScore)
			
			assert.Equal(t, models.EventTypeOpponentGoal, u.Event[0].Type)
			assert.Equal(t, int64(1), u.Event[0].TeamScore)
			assert.Equal(t, int64(0), u.Event[0].OpponentScore)
		})

		t.Run("17 start point offense", func(t *testing.T) {
			err = event.StartingLine(gameCtx, playerSlugs...)
			assert.Nil(t, err)
			if err != nil {
				return
				}
			<-sub.EventChan
			err = event.StartPointOffense(gameCtx)
			assert.Nil(t, err)
			if err != nil {
				return
			}
			u := <-sub.EventChan
			assert.Equal(t, 1, len(u.Event))
			assert.Equal(t, u.Id, u.Game.LastEvent.String)
			assert.Equal(t, models.GameLiveStatusTeamPossession, u.Game.LiveStatus)
			assert.Equal(t, models.EventTypePointStart, u.Event[0].Type)
			assert.Equal(t, int64(1), u.Event[0].TeamScore)
			assert.Equal(t, int64(1), u.Event[0].OpponentScore)
		})

		t.Run("18 team goal", func(t *testing.T) {
			err = event.Goal(gameCtx, playerSlugs[1], playerSlugs[0])
			assert.Nil(t, err)
			if err != nil {
				return
			}
			u := <-sub.EventChan
			assert.Equal(t, 2, len(u.Event))
			assert.Equal(t, u.Id, u.Game.LastEvent.String)
			assert.Equal(t, models.GameLiveStatusPrePoint, u.Game.LiveStatus)
			assert.Equal(t, int64(2), u.Game.TeamScore)
			assert.Equal(t, int64(1), u.Game.OpponentScore)

			assert.Equal(t, models.EventTypeGoal, u.Event[0].Type)
			assert.Equal(t, playerIds[1], u.Event[0].Player.Int64)
			assert.Equal(t, int64(1), u.Event[0].TeamScore)
			assert.Equal(t, int64(1), u.Event[0].OpponentScore)
			assert.Equal(t, u.Id, u.Event[0].Batch.String)

			assert.Equal(t, models.EventTypeAssist, u.Event[1].Type)
			assert.Equal(t, playerIds[0], u.Event[1].Player.Int64)
			assert.Equal(t, int64(1), u.Event[0].TeamScore)
			assert.Equal(t, int64(1), u.Event[0].OpponentScore)
			assert.Equal(t, u.Id, u.Event[1].Batch.String)
		})
	})
}
