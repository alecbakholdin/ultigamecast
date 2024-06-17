package ctxvar

import (
	"context"
	"testing"
	"ultigamecast/internal/models"

	"github.com/stretchr/testify/assert"
)

func TestUrlStandardValues(t *testing.T) {
	ctx := context.WithValue(context.Background(), Team, &models.Team{Slug: "team-slug", ID: 1})
	ctx = context.WithValue(ctx, Tournament, &models.Tournament{Slug: "tournament-slug", ID: 2})
	ctx = context.WithValue(ctx, Player, &models.Player{Slug: "player-slug", ID: 3})
	ctx = context.WithValue(ctx, Game, &models.Game{Slug: "game-slug", ID: 4})
	ctx = context.WithValue(ctx, Event, &models.Event{ID: "random-string"})

	assert.Equal(t, "/teams/team-slug/schedule/tournament-slug/players/player-slug/games/game-slug/events/random-string/url", Url(ctx, Team, Tournament, Player, Game, Event, "url"))
	assert.Equal(t, "/", Url(ctx, ""))
	assert.Equal(t, "/teams/team-slug", Url(ctx, "", Team))
	assert.Panics(t, func() { Url(ctx, "", HttpMethod) })
	assert.Panics(t, func() { Url(context.Background(), "", Team) })
}

func TestUrlUnsupportedValuePanics(t *testing.T) {
	ctx := context.WithValue(context.Background(), Team, &models.Team{Slug: "team-slug", ID: 1})
	assert.Panics(t, func() { Url(ctx, "", HttpMethod) })
}

func TestUrlMissingValuePanics(t *testing.T) {
	assert.Panics(t, func() { Url(context.Background(), "", Team) })
}
