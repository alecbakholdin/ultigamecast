package repository

import (
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
}

func NewGame(app core.App) *Game {
	return &Game{
		app:        app,
		dao:        app.Dao(),
		collection: mustGetCollection(app.Dao(), "games"),
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
	game.SetIsCompleted(gameDto.GamesIsCompleted)

	if err := g.dao.SaveRecord(game.Record); err != nil {
		return nil, err
	}
	return game, nil
}

func (g *Game) Update(id string, gameDto *dto.Games) (game *modelspb.Games, err error) {
	if record, err := g.dao.FindRecordById(g.collection.Id, id); err != nil {
		return nil, err
	} else {
		game = toGame(record)
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
	game.SetIsCompleted(gameDto.GamesIsCompleted)

	if err := g.dao.SaveRecord(game.Record); err != nil {
		return nil, err
	}
	return game, nil
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
