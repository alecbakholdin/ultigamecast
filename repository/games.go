package repository

import (
	"strings"
	"ultigamecast/pbmodels"

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
	dao := app.Dao()
	collection := mustGetCollection(app.Dao(), "games")

	return &Game{
		app:        app,
		dao:        dao,
		collection: collection,
	}
}

func (g *Game) GetOneById(id string) (*pbmodels.Games, error) {
	game := &pbmodels.Games{}
	err := g.gameQuery().Where(dbx.HashExp{"id": id}).Limit(1).One(game)
	if err != nil {
		return nil, err
	}
	return game, err
}

func (g *Game) GetAllByTeamSlug(teamSlug string) ([]*pbmodels.Games, error) {
	q := g.gameQuery()
	q.InnerJoin("tournaments", dbx.NewExp("tournaments.id = games.tournament"))
	q.InnerJoin("teams", dbx.NewExp("teams.id = tournaments.team"))
	q.Where(dbx.HashExp{"teams.slug": strings.ToLower(teamSlug)})
	q.OrderBy("start DESC")

	games := make([]*pbmodels.Games, 0)
	if err := q.All(&games); err != nil {
		return nil, err
	}
	return games, nil
}

func (g *Game) Create(game *pbmodels.Games) error {
	if err := g.dao.DB().Model(game).Exclude("Id").Insert(); err != nil {
		return err
	}
	return nil
}

func (g *Game) Update(game *pbmodels.Games) (err error) {
	if err := g.dao.DB().Model(game).Exclude("Id", "Tournament").Update(); err != nil {
		return err
	}
	return nil
}

func (g *Game) GetAllByTeamAndTournamentSlugs(teamSlug string, tournamentSlug string) ([]*pbmodels.Games, error) {
	games := make([]*pbmodels.Games, 0)
	err := g.gameQuery().InnerJoin(
		"tournaments", dbx.NewExp("tournaments.id = games.tournament"),
	).InnerJoin(
		"teams", dbx.NewExp("tournaments.team = teams.id"),
	).Where(dbx.HashExp{
		"teams.slug":       strings.ToLower(teamSlug),
		"tournaments.slug": strings.ToLower(tournamentSlug),
	}).All(&games)
	if err != nil {
		return nil, err
	}
	return games, nil
}

func (g *Game) gameQuery() *dbx.SelectQuery {
	return g.dao.ModelQuery(&pbmodels.Games{})
}
