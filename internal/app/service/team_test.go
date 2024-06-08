package service

import (
	"errors"
	"testing"
	"ultigamecast/test/testctx"
	"ultigamecast/test/testdb"

	"github.com/stretchr/testify/assert"
)

func TestTeam(t *testing.T) {
	te := NewTeam(testdb.DB())
	ctx := testctx.LoadUser(te.q)

	t.Run("create team", func(t *testing.T) {
		team, err := te.CreateTeam(ctx, "Test Create Team", "organization")
		assert.Nil(t, err, "error creating team: %s", err)
		assert.Equal(t, "Test Create Team", team.Name)
		assert.Equal(t, "test-create-team", team.Slug)
		assert.Equal(t, "organization", team.Organization.String)

		_, err = te.CreateTeam(ctx, "test create team", "noorg")
		assert.True(t, errors.Is(err, ErrTeamExists), "incorrect error when creating duplicate team: %s", err)
	})

	t.Run("list teams", func(t *testing.T) {
		teams, err := te.GetTeams(ctx)
		assert.Nil(t, err, "error getting teams: %s", err)
		_, err = te.CreateTeam(ctx, "TestListTeams", "or")
		assert.Nil(t, err, "error creating team: %s", err)
		teams2, err := te.GetTeams(ctx)
		assert.Nil(t, err, "error getting teams2: %s", err)
		assert.Equal(t, len(teams)+1, len(teams2))
	})

	t.Run("update name", func(t *testing.T) {
		team, err := te.CreateTeam(ctx, "new team name", "orgo")
		assert.Nil(t, err, "error creating team: %s", err)

		teamCtx := testctx.Load(ctx, team)
		updatedTeam, err := te.UpdateName(teamCtx, "New Name")
		assert.Nil(t, err, "error updating team name: %s", err)
		assert.Equal(t, "new-name", updatedTeam.Slug)
		assert.Equal(t, "New Name", updatedTeam.Name)
	})

	t.Run("update organization", func(t *testing.T) {
		team, err := te.CreateTeam(ctx, "new team name", "orgo")
		assert.Nil(t, err, "error creating team: %s", err)

		teamCtx := testctx.Load(ctx, team)
		updatedTeam, err := te.UpdateOrganization(teamCtx, "organization")
		assert.Nil(t, err, "error updating team organization: %s", err)
		assert.Equal(t, "organization", updatedTeam.Organization.String)
	})
}
