package handlers

import (
	"fmt"
	"net/http"
	"ultigamecast/modelspb"
	"ultigamecast/modelspb/dto"
	"ultigamecast/repository"
	"ultigamecast/validation"
	view "ultigamecast/view/team/tournaments/games"

	"github.com/labstack/echo/v5"
)

type Games struct {
	teamRepo      *repository.Team
	touramentRepo *repository.Tournament
	gameRepo      *repository.Game
}

func NewGames(te *repository.Team, to *repository.Tournament, g *repository.Game) *Games {
	return &Games{
		teamRepo:      te,
		touramentRepo: to,
		gameRepo:      g,
	}
}

func (g *Games) Routes(tournamentGroup *echo.Group) *echo.Group {
	tournamentGroup.GET("/newGame", g.getNewGameModal)
	tournamentGroup.GET("/games", g.getGames)
	tournamentGroup.POST("/games", g.createGame)

	gameGroup := tournamentGroup.Group("/games/:gameId")
	gameGroup.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.PathParam("gameId") == "" {
				return echo.NewHTTPError(http.StatusBadRequest, "missing game in request")
			}
			return next(c)
		}
	})
	gameGroup.GET("/edit", g.getEditGameModal);
	gameGroup.PUT("", g.updateGame);
	return gameGroup
}

func (g *Games) getGames(c echo.Context) (err error) {
	teamSlug := c.PathParam("teamSlug")
	tournamentSlug := c.PathParam("tournamentSlug")

	if games, err := g.gameRepo.GetAllByTeamAndTournamentSlugs(teamSlug, tournamentSlug); err != nil {
		return echo.NewHTTPErrorWithInternal(http.StatusInternalServerError, err, "could not get tournament games")
	} else {
		return view.TournamentGameList(teamSlug, tournamentSlug, games).Render(c.Request().Context(), c.Response().Writer)
	}
}

func (g *Games) getNewGameModal(c echo.Context) (err error) {
	TriggerOpenModal(c)
	return view.CreateEditGameDialogContent(c, true, dto.Games{}).Render(c.Request().Context(), c.Response().Writer)
}

func (g *Games) getEditGameModal(c echo.Context) (err error) {
	gameId := c.PathParam("gameId")

	if game, err := g.gameRepo.GetOneById(gameId); err != nil {
		return echo.NewHTTPErrorWithInternal(http.StatusInternalServerError, err, "unexpected error")
	} else {
		TriggerOpenModal(c)
		return view.CreateEditGameDialogContent(c, false, *dto.DtoFromGame(c, game)).Render(c.Request().Context(), c.Response().Writer)
	}
}

func (g *Games) createGame(c echo.Context) (err error) {
	var (
		payload    dto.Games
		tournament *modelspb.Tournaments
	)
	if err = dto.BindGameDto(c, &payload); err != nil {
		c.Echo().Logger.Error(fmt.Errorf("error binding game dto: %s", err))
		return err
	}

	if tournament, err = g.touramentRepo.GetOneBySlug(payload.TeamSlug, payload.TournamentSlug); err != nil {
		c.Echo().Logger.Error(fmt.Errorf("error finding tournament [%s, %s]: %s", payload.TeamSlug, payload.TournamentSlug, err))
		validation.AddFormErrorString(c, "could not find associated tournament")
	}

	if !validation.IsFormValid(c) {
		return view.GameForm(c, true, payload).Render(c.Request().Context(), c.Response().Writer)
	} else if _, err := g.gameRepo.Create(tournament, &payload); err != nil {
		c.Echo().Logger.Error(fmt.Errorf("error creating game: %s", err))
		validation.AddFormErrorString(c, "could not create game")
	} else {
		TriggerCloseModal(c)
		triggerGameListRefresh(c, payload.TournamentSlug)
	}

	return view.GameForm(c, true, payload).Render(c.Request().Context(), c.Response().Writer)
}

func (g *Games) updateGame(c echo.Context) (err error){ 
	var (
		payload    dto.Games
	)
	if err = dto.BindGameDto(c, &payload); err != nil {
		c.Echo().Logger.Error(fmt.Errorf("error binding game dto: %s", err))
		return err
	}

	if !validation.IsFormValid(c) {
		return view.GameForm(c, false, payload).Render(c.Request().Context(), c.Response().Writer)
	} else if _, err := g.gameRepo.Update(payload.GameID, &payload); err != nil {
		c.Echo().Logger.Error(fmt.Errorf("error updating game: %s", err))
		validation.AddFormErrorString(c, "unexpected error updating game")
	} else {
		TriggerCloseModal(c)
		triggerGameListRefresh(c, payload.TournamentSlug)
	}

	return view.GameForm(c, true, payload).Render(c.Request().Context(), c.Response().Writer)
}

func triggerGameListRefresh(c echo.Context, tournamentSlug string) {
	c.Response().Header().Add("HX-Trigger", fmt.Sprintf("refreshgames-%s", tournamentSlug))
}