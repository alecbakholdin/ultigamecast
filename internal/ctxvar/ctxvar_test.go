package ctxvar

import (
	"context"
	"testing"
	"ultigamecast/internal/models"

	"github.com/stretchr/testify/assert"
)

func TestUrl(t *testing.T) {
	t.Run("standard values", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), Team, &models.Team{Slug: "team-slug", ID: 1})
		ctx = context.WithValue(ctx, Tournament, &models.Tournament{Slug: "tournament-slug", ID: 2})
		ctx = context.WithValue(ctx, Player, &models.Player{Slug: "player-slug", ID: 3})
		ctx = context.WithValue(ctx, Game, &models.Game{Slug: "game-slug", ID: 4})
		ctx = context.WithValue(ctx, Event, &models.Event{ID: "random-string"})

		assert.Equal(t, "/teams/team-slug/schedule/tournament-slug/players/player-slug/schedule/game-slug/events/random-string/url", Url(ctx, Team, Tournament, Player, Game, Event, "url"))
		assert.Equal(t, "/", Url(ctx, ""))
		assert.Equal(t, "/teams/team-slug", Url(ctx, "", Team))
		assert.Panics(t, func() { Url(ctx, "", HttpMethod) })
		assert.Panics(t, func() { Url(context.Background(), "", Team) })
	})

	t.Run("prefix strips slashes", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), Team, &models.Team{Slug: "team-slug", ID: 1})
		assert.Equal(t, "/prefix/teams/team-slug", Url(ctx, "/prefix/", Team))
		assert.Equal(t, "/teams/team-slug/teams/team-slug", Url(ctx, Url(ctx, Team), Team))
	})

	t.Run("unsupported value panics", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), Team, &models.Team{Slug: "team-slug", ID: 1})
		assert.Panics(t, func() { Url(ctx, "", HttpMethod) })
	})

	t.Run("missing value panics", func(t *testing.T) {
		assert.Panics(t, func() { Url(context.Background(), "", Team) })
	})
}
