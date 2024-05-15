package handlers

import (
	"fmt"
	"net/http"
	"ultigamecast/pbmodels"
	"ultigamecast/repository"
	"ultigamecast/service"
	"ultigamecast/validation"
	view "ultigamecast/view/team/tournaments/games"
	"ultigamecast/view/team/tournaments/games/game"

	"github.com/labstack/echo/v5"
)

type Games struct {
	teamRepo         *repository.Team
	touramentRepo    *repository.Tournament
	gameService      *service.Games
	eventsService    *service.Events
	liveGamesService *service.LiveGames
}

func NewGames(te *repository.Team, to *repository.Tournament, g *service.Games, e *service.Events, l *service.LiveGames) *Games {
	return &Games{
		teamRepo:         te,
		touramentRepo:    to,
		gameService:      g,
		eventsService:    e,
		liveGamesService: l,
	}
}

func (g *Games) Routes(tournamentGroup *echo.Group) *echo.Group {
	tournamentGroup.GET("/newGame", g.getNewGameModal)
	tournamentGroup.GET("/games", g.getGames)
	tournamentGroup.POST("/games", g.createGame)

	gameGroup := tournamentGroup.Group("/games/:gamesId")
	gameGroup.GET("/edit", g.getEditGameModal)
	gameGroup.PUT("", g.updateGame)
	gameGroup.GET("", g.getGame)

	gameGroup.GET("/sse", g.gameSse)
	gameGroup.POST("/opponentScore", g.opponentScore)
	return gameGroup
}

func (g *Games) getGames(c echo.Context) (err error) {
	teamSlug := c.PathParam("teamsSlug")
	tournamentSlug := c.PathParam("tournamentsSlug")

	if games, err := g.gameService.GetAllBySlugs(teamSlug, tournamentSlug); err != nil {
		return echo.NewHTTPErrorWithInternal(http.StatusInternalServerError, err, "could not get tournament games")
	} else {
		return view.TournamentGameList(teamSlug, tournamentSlug, games).Render(c.Request().Context(), c.Response().Writer)
	}
}

func (g *Games) getNewGameModal(c echo.Context) (err error) {
	TriggerOpenModal(c)
	return view.CreateEditGameDialogContent(c, true, &pbmodels.Games{}).Render(c.Request().Context(), c.Response().Writer)
}

func (g *Games) getEditGameModal(c echo.Context) (err error) {
	gameId := c.PathParam("gamesId")

	if game, err := g.gameService.GetOneById(gameId); err != nil {
		return echo.NewHTTPErrorWithInternal(http.StatusInternalServerError, err, "unexpected error")
	} else {
		TriggerOpenModal(c)
		return view.CreateEditGameDialogContent(c, false, game).Render(c.Request().Context(), c.Response().Writer)
	}
}

func (g *Games) createGame(c echo.Context) (err error) {
	var (
		game           pbmodels.Games
		teamSlug       = c.PathParam("teamsSlug")
		tournamentSlug = c.PathParam("tournamentsSlug")
	)
	if err = c.Bind(&game); err != nil {
		c.Echo().Logger.Error(fmt.Errorf("error binding game dto: %s", err))
		return err
	}
	g.gameService.Validate(c, &game)

	if !validation.IsFormValid(c) {
		return view.GameForm(c, true, &game).Render(c.Request().Context(), c.Response().Writer)
	} else if err := g.gameService.Create(teamSlug, tournamentSlug, &game); err != nil {
		c.Echo().Logger.Error(fmt.Errorf("error creating game: %s", err))
		validation.AddFormErrorString(c, "unexpected error creating game")
	} else {
		TriggerCloseModal(c)
		triggerGameListRefresh(c, tournamentSlug)
	}

	return view.GameForm(c, true, &game).Render(c.Request().Context(), c.Response().Writer)
}

func (g *Games) updateGame(c echo.Context) (err error) {
	var (
		game           pbmodels.Games
		gameId         = c.PathParam("gamesId")
		tournamentSlug = c.PathParam("tournamentsSlug")
	)
	if err = c.Bind(&game); err != nil {
		c.Echo().Logger.Error(fmt.Errorf("error binding game dto: %s", err))
		return err
	}
	game.Id = gameId

	g.gameService.Validate(c, &game)

	if !validation.IsFormValid(c) {
		return view.GameForm(c, false, &game).Render(c.Request().Context(), c.Response().Writer)
	} else if err := g.gameService.Update(game.Id, &game); err != nil {
		c.Echo().Logger.Error(fmt.Errorf("error updating game: %s", err))
		validation.AddFormErrorString(c, "unexpected error updating game")
	} else {
		TriggerCloseModal(c)
		triggerGameListRefresh(c, tournamentSlug)
	}

	if validation.IsFormValid(c) && game.Status != pbmodels.GamesStatusLive {
		g.liveGamesService.RemoveLiveGame(gameId)
	}

	return view.GameForm(c, false, &game).Render(c.Request().Context(), c.Response().Writer)
}

func triggerGameListRefresh(c echo.Context, tournamentSlug string) {
	c.Response().Header().Add("HX-Trigger", fmt.Sprintf("refreshgames-%s", tournamentSlug))
}

func (g *Games) getGame(c echo.Context) (err error) {
	var (
		gameModel *pbmodels.Games
		events    []*pbmodels.Events
		gameId = c.PathParam("gamesId")
	)
	if gameModel, err = g.gameService.GetOneById(gameId); err != nil {
		return echo.NewHTTPErrorWithInternal(http.StatusInternalServerError, err, "Unexpected error")
	}
	if events, err = g.eventsService.GetAllByGame(gameId); err != nil {
		return echo.NewHTTPErrorWithInternal(http.StatusInternalServerError, err, "Unexpected error")
	}

	return game.GameRoot(c, gameModel, events).Render(c.Request().Context(), c.Response())
}

func (g *Games) gameSse(c echo.Context) (err error) {
	sub, err := g.liveGamesService.Subscribe(c.PathParam("gamesId"))
	if err != nil {
		return echo.NewHTTPErrorWithInternal(http.StatusInternalServerError, err, "Unexpected error")
	}
	defer g.liveGamesService.Unsubscribe(sub)

	w := c.Response()
	w.Header().Add("Content-Type", "text/event-stream")
	w.Header().Add("Cache-Control", "no-cache")
	w.Header().Add("Connection", "keep-alive")

	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "GET")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")

outer_sse:
	for {
		select {
		case e := <-sub.EventUpdates:
			w.Write([]byte(fmt.Sprintf("event: %s\ndata: ", e.Event.Id)))
			game.GameEvent(e.Event).Render(c.Request().Context(), w)
			w.Write([]byte("\n\n"))
		case g := <-sub.GameUpdates:
			for _, f := range g.Fields {
				w.Write([]byte(fmt.Sprintf("event: %s\ndata: ", f)))
				game.GameUpdate(f, g.Game).Render(c.Request().Context(), w)
				w.Write([]byte("\n\n"))
			}
			w.Flush()
		case <-sub.Close:
			c.Echo().Logger.Write([]byte("done"))
			break outer_sse
		}
	}

	return nil
}

func (g *Games) opponentScore(c echo.Context) (err error) {
	var (
		gameId  = c.PathParam("gamesId")
		message = c.QueryParam("message")
	)

	err = g.liveGamesService.OpponentScored(gameId, message)
	if err != nil {
		return echo.NewHTTPErrorWithInternal(http.StatusInternalServerError, err, "unexpected error")
	}

	return nil
}
