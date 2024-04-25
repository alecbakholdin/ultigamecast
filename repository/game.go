package repository

import (
	"log"
	"strings"
	"ultigamecast/modelspb"
	"ultigamecast/modelspb/dto"

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

func (g *Game) Create(tournament *modelspb.Tournaments, gameDto *dto.Games) (*modelspb.Games, error) {
	game := toGame(models.NewRecord(g.collection))

	game.SetTournament(tournament.Record.GetId())
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

func (g *Game) GetAllByTeamAndTournamentSlugs(teamSlug string, tournamentSlug string) ([]*modelspb.Games, error) {
	records, err := g.dao.FindRecordsByFilter(
		g.collection.Id,
		"tournament.slug = {:tournamentSlug} && tournament.team.slug = {:teamSlug}",
		"+start_time",
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
