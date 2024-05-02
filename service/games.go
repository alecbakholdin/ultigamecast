package service

import (
	"fmt"
	"ultigamecast/pbmodels"
	"ultigamecast/repository"
	"ultigamecast/validation"

	"github.com/labstack/echo/v5"
)

type Games struct {
	TournamentRepo *repository.Tournament
	GameRepo       *repository.Game
}

func NewGames(to *repository.Tournament, g *repository.Game) *Games {
	return &Games{
		TournamentRepo: to,
		GameRepo:       g,
	}
}

func (g *Games) GetOneById(id string) (*pbmodels.Games, error) {
	return g.GameRepo.GetOneById(id)
}

func (g *Games) GetAllBySlugs(teamSlug, tournamntSlug string) ([]*pbmodels.Games, error) {
	return g.GameRepo.GetAllByTeamAndTournamentSlugs(teamSlug, tournamntSlug)
}

func (g *Games) Create(teamSlug, tournamentSlug string, game *pbmodels.Games) (err error) {
	var tournament *pbmodels.Tournaments
	if tournament, err = g.TournamentRepo.GetOneBySlug(teamSlug, tournamentSlug); err != nil {
		return fmt.Errorf("error finding tournament [%s] [%s]: %s", teamSlug, tournamentSlug, err)
	}
	game.Tournament = tournament.Id
	if err = g.GameRepo.Create(game); err != nil {
		return fmt.Errorf("error creating game for tournament [%s] [%s]: %s", teamSlug, tournamentSlug, err)
	}
	return nil
}

func (g *Games) Update(gameId string, game *pbmodels.Games) (err error) {
	if gameId == "" {
		return fmt.Errorf("gameId is empty")
	}

	game.Id = gameId
	if err = g.GameRepo.Update(game); err != nil {
		return fmt.Errorf("error updating game [%s]: %s", gameId, err)
	}
	return nil
}

func (g *Games) Validate(c echo.Context, game *pbmodels.Games) {
	var err error
	if game.StartTime, err = game.GetStartTimeDt(); err != nil {
		c.Echo().Logger.Error(fmt.Errorf("error parsing start datetime: %s", err))
		validation.AddFieldErrorString(c, "start_time", "invalid format")
	}
}
