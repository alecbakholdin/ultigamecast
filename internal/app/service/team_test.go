package service

import (
	"errors"
	"testing"
	"ultigamecast/test/testctx"
	"ultigamecast/test/testdb"

	"github.com/stretchr/testify/assert"
)

func TestCreateTeam(t *testing.T) {
	te := NewTeam(testdb.DB())
	ctx := testctx.LoadUser(te.q)
	team, err := te.CreateTeam(ctx, "Test Create Team", "organization")
	assert.Nil(t, err, "error creating team: %s", err)
	assert.Equal(t, "Test Create Team", team.Name)
	assert.Equal(t, "test-create-team", team.Slug)
	assert.Equal(t, "organization", team.Organization.String)

	_, err = te.CreateTeam(ctx, "test create team", "noorg")
	assert.True(t, errors.Is(err, ErrTeamExists), "incorrect error when creating duplicate team: %s", err)
}

func TestListTeams(t *testing.T) {
	te := NewTeam(testdb.DB())
	ctx := testctx.LoadUser(te.q)
	teams, err := te.GetTeams(ctx)
	assert.Nil(t, err, "error getting teams: %s", err)
	_, err = te.CreateTeam(ctx, "TestListTeams", "or")
	assert.Nil(t, err, "error creating team: %s", err)
	teams2, err := te.GetTeams(ctx)
	assert.Nil(t, err, "error getting teams2: %s", err)
	assert.Equal(t, len(teams) + 1, len(teams2))
}