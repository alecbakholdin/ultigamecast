package repository

import (
	"log"
	"strings"
	"ultigamecast/modelspb"
	"ultigamecast/modelspb/dto"
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
	liveGame   *LiveGame
	allEvents  []string
}

func NewGame(app core.App, l *LiveGame) *Game {
	dao := app.Dao()
	collection := mustGetCollection(app.Dao(), "games")

	events := make([]string, len(collection.Schema.AsMap()))
	i := 0
	for key := range collection.Schema.AsMap() {
		events[i] = key
		i++
	}
	log.Printf("%v\n", events)

	return &Game{
		app:        app,
		dao:        dao,
		collection: collection,
		liveGame:   l,
		allEvents:  events,
	}
}

func (g *Game) GetOneById(id string) (*modelspb.Games, error) {
	if record, err := g.dao.FindRecordById(g.collection.Id, id); err != nil {
		return nil, err
	} else {
		return toGame(record), nil
	}
}

func (g *Game) GetAllForTeamBySlug(teamSlug string) ([]*pbmodels.Games, error) {
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

func (g *Game) Create(tournament *pbmodels.Tournaments, gameDto *dto.Games) (*modelspb.Games, error) {

	game := toGame(models.NewRecord(g.collection))

	game.SetTournament(tournament.Id)
	game.SetOpponent(gameDto.GameOpponent)
	game.SetTeamScore(gameDto.GameTeamScore)
	game.SetOpponentScore(gameDto.GameOpponentScore)
	game.SetHalfCap(gameDto.GameHalfCap)
	game.SetSoftCap(gameDto.GameSoftCap)
	game.SetHardCap(gameDto.GameHardCap)
	game.SetWindMph(gameDto.GameWindMph)
	game.SetTempF(gameDto.GameTempF)
	game.SetStartTime(gameDto.GameStartTimeDt)
	game.SetStatus(gameDto.GamesStatus)

	if err := g.dao.SaveRecord(game.Record); err != nil {
		return nil, err
	}
	g.liveGame.trigger(g.allEvents, game)
	return game, nil
}

func (g *Game) Update(id string, gameDto *dto.Games) (game *modelspb.Games, err error) {
	if game, err = g.GetOneById(id); err != nil {
		return nil, err
	}

	game.SetOpponent(gameDto.GameOpponent)
	game.SetTeamScore(gameDto.GameTeamScore)
	game.SetOpponentScore(gameDto.GameOpponentScore)
	game.SetHalfCap(gameDto.GameHalfCap)
	game.SetSoftCap(gameDto.GameSoftCap)
	game.SetHardCap(gameDto.GameHardCap)
	game.SetWindMph(gameDto.GameWindMph)
	game.SetTempF(gameDto.GameTempF)
	game.SetStartTime(gameDto.GameStartTimeDt)
	game.SetStatus(gameDto.GamesStatus)

	if err := g.dao.SaveRecord(game.Record); err != nil {
		return nil, err
	}
	g.liveGame.trigger(g.allEvents, game)
	return game, nil
}

func (g *Game) UpdateField(id string, field string, value any) (err error) {
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

func toGame(record *models.Record) *modelspb.Games {
	return &modelspb.Games{
		Record: record,
	}
}
