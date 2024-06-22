package service

import (
	"testing"
	"ultigamecast/test/testctx"
	"ultigamecast/test/testdb"

	"github.com/stretchr/testify/assert"
)

func TestCreatePlayer(t *testing.T) {
	p := NewPlayer(testdb.DB())
	ctx := testctx.LoadTeam(p.q)
	players1, err := p.GetTeamPlayers(ctx)
	assert.Nil(t, err, "error getting team players: %s", err)
	assert.EqualValues(t, 0, len(players1))

	player, err := p.CreatePlayer(ctx, "Player name")
	assert.Nil(t, err, "error creating player: %s", err)
	assert.EqualValues(t, "Player name", player.Name)
	assert.EqualValues(t, "player-name", player.Slug)
	assert.EqualValues(t, 0, player.Order)

	player2, err := p.CreatePlayer(ctx, "player name")
	assert.Nil(t, err, "error creating player2: %s", err)
	assert.EqualValues(t, "player name", player2.Name)
	assert.EqualValues(t, "player-name-2", player2.Slug)
	assert.EqualValues(t, 1, player2.Order)

	players2, err := p.GetTeamPlayers(ctx)
	assert.Nil(t, err)
	assert.EqualValues(t, 2, len(players2))
	assert.EqualValues(t, "player-name", players2[0].Slug)
	assert.EqualValues(t, "player-name-2", players2[1].Slug)
}


func TestUpdatePlayerOrder(t *testing.T) {
	p := NewPlayer(testdb.DB())
	ctx := testctx.LoadTeam(p.q)

	players1, err := p.GetTeamPlayers(ctx)
	assert.Nil(t, err, "error getting team players: %s")
	players1Ids := make([]int64, len(players1))
	for i, p := range players1 {
		players1Ids[(i + 1) % len(players1)] = p.ID
	}

	err = p.UpdatePlayerOrder(ctx, players1Ids)
	assert.Nil(t, err, "error updating player order: %s", err)

	players2, err := p.GetTeamPlayers(ctx)
	assert.Nil(t, err, "error getting team players2: %s")

	players2Ids := make([]int64, len(players2))
	for i, p := range players2 {
		players2Ids[i] = p.ID
	}

	assert.Equal(t, players1Ids, players2Ids)
}