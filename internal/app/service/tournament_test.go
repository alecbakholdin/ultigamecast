package service

import (
	"fmt"
	"testing"
	"time"
	"ultigamecast/test/testctx"
	"ultigamecast/test/testdb"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func TestTournament(t *testing.T) {
	to := NewTournament(testdb.DB())
	t.Run("create simple", func(t *testing.T) {
		ctx := testctx.LoadTeam(to.q)
		tournament, err := to.CreateTournament(ctx, "Random name", "Jan 2, 2024 - Jan 3, 2024")
		assert.Nil(t, err, "error creating tournament")
		assert.Equal(t, "Random name", tournament.Name)
		assert.Equal(t, "random-name", tournament.Slug)
		assert.Equal(t, mustParseTime("Jan 2, 2006", "Jan 2, 2024"), tournament.StartDate.Time)
		assert.Equal(t, mustParseTime("Jan 2, 2006", "Jan 3, 2024"), tournament.EndDate.Time)
	})
	t.Run("create no time", func(t *testing.T) {
		ctx := testctx.LoadTeam(to.q)
		tournament, err := to.CreateTournament(ctx, "name", "")
		assert.Nil(t, err, "error creating tournament")
		assert.Equal(t, "name", tournament.Name)
		assert.Equal(t, "name", tournament.Slug)
		assert.False(t, tournament.StartDate.Valid, "start date is valid when it shouldn't be")
		assert.False(t, tournament.EndDate.Valid, "end date is valid when it shouldn't be")
	})
	t.Run("create duplicate names", func(t *testing.T) {
		ctx := testctx.LoadTeam(to.q)
		t1, err := to.CreateTournament(ctx, "Dupe Name", "")
		assert.Nil(t, err, "error creating tournament")
		assert.Equal(t, "dupe-name", t1.Slug)
		t2, err := to.CreateTournament(ctx, "dupe name", "")
		assert.Nil(t, err, "error creating tournament")
		assert.Equal(t, "dupe-name-2", t2.Slug)
	})
	t.Run("datum creation", func(t *testing.T) {
		ctx := testctx.LoadTournament(to.q)
		_, err := to.Data(ctx)
		assert.Nil(t, err, "error fetching data") // doesnt return error when empty
		datum, err := to.NewDatum(ctx)
		assert.Nil(t, err, "error creating datum")
		data, err := to.Data(ctx)
		assert.Nil(t, err, "error fetching data")
		assert.Equal(t, 1, len(data))
		assert.Equal(t, datum.ID, data[0].ID)
	})
	t.Run("update data order", func(t *testing.T) {
		ctx := testctx.LoadTournament(to.q)
		var ids []int64
		for range 5 {
			datum, err := to.NewDatum(ctx)
			assert.Nil(t, err, "error creating datum")
			ids = append(ids, datum.ID)
		}

		dataPreOrder, err := to.Data(ctx)
		assert.Nil(t, err, "error fetching data pre-order")
		idsPreOrder := make([]int64, len(dataPreOrder))
		for i, d := range dataPreOrder {
			idsPreOrder[i] = d.ID
		}
		assert.Equal(t, ids, idsPreOrder)

		temp := ids[1]
		ids[1] = ids[2]
		ids[2] = temp
		assert.Nil(t, to.UpdateDataOrder(ctx, ids), "error updating order")

		dataPostOrder, err := to.Data(ctx)
		assert.Nil(t, err, "error fetching data post-order")
		idsPostOrder := make([]int64, len(dataPostOrder))
		for i, d := range dataPostOrder {
			idsPostOrder[i] = d.ID
		}
		assert.Equal(t, ids, idsPostOrder)
	})
	t.Run("get schedule", func(t *testing.T) {
		ctx := testctx.LoadTeam(to.q)
		tournament, err := to.CreateTournament(ctx, "get schedule", "Jan 1, 2024 - Jan 2, 2024")
		assert.Nil(t, err, "error creating tournament")
		withTournament := testctx.Load(ctx, tournament)
		g := NewGame(testdb.DB())
		game, err := g.CreateGame(withTournament, "opp", "2024-01-02T15:04", "America/New_York", 1, 2, 3)
		assert.Nil(t, err, "error creating game")
		schedule, err := to.GetSchedule(ctx)
		assert.Nil(t, err, "error getting schedule")
		assert.Equal(t, 1, len(schedule))
		assert.Equal(t, *tournament.Tournament, *schedule[0].Tournament)
		assert.Equal(t, 1, len(schedule[0].Games))
		assert.Equal(t, *game, schedule[0].Games[0])
	})
}

func mustParseTime(layout, val string) time.Time {
	time, err := time.Parse(layout, val)
	if err != nil {
		panic(fmt.Errorf("error parsing time: %w", err))
	}
	return time
}

