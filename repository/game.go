package repository

import (
	"strings"
	"ultigamecast/modelspb"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
)

type Game struct {
	app        core.App
	dao        *daos.Dao
	collection *models.Collection
}

func NewGame(app core.App) *Game {
	return &Game{
		app:        app,
		dao:        app.Dao(),
		collection: mustGetCollection(app.Dao(), "games"),
	}
}

func (g *Game) GetAllByTeamAndTournamentSlugs(teamSlug string, tournamentSlug string) ([]*modelspb.Games, error) {
	records, err := g.dao.FindRecordsByFilter(
		g.collection.Id,
		"tournament.slug = {:tournamentSlug} && tournament.team.slug = {:teamSlug}",
		"-start",
		0,
		0,
		dbx.Params{
			"teamSlug":       strings.ToLower(teamSlug),
			"tournamentSlug": strings.ToLower(tournamentSlug),
		},
	)
	if err != nil {
		return []*modelspb.Games{}, err
	}
	return toArr(records, toGame), nil
}

func toGame(record *models.Record) *modelspb.Games {
	return &modelspb.Games{
		Record: record,
	}
}
