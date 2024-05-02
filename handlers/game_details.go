package handlers

import (
	"fmt"
	"log"
	"ultigamecast/modelspb/dto"
	"ultigamecast/repository"

	gameview "ultigamecast/view/team/tournaments/games/game"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v5"
)

type GameDetails struct {
	gameRepo *repository.Game
	liveGame *repository.LiveGame
}

func NewGameDetails(g *repository.Game, l *repository.LiveGame) *GameDetails {
	return &GameDetails{
		gameRepo: g,
		liveGame: l,
	}
}

func (g *GameDetails) Routes(gameGroup *echo.Group) *echo.Group {
	gameGroup.GET("", g.getGame)
	gameGroup.GET("/sse", g.sseGame)

	return gameGroup
}

func (g *GameDetails) getGame(c echo.Context) (err error) {
	return nil
}

func (g *GameDetails) sseGame(c echo.Context) (err error) {
	log.Printf("SSE client connected, ip: %v", c.RealIP())

	w := c.Response()
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	gameId := c.PathParam("gameId")

	subId, events := g.liveGame.Subscribe(c, gameId)
	defer g.liveGame.Unsubscribe(subId)
	for {
		select {
		case <-c.Request().Context().Done():
			log.Printf("SSE client disconnected, ip: %v", c.RealIP())
			return nil
		case event := <-events:
			if com := getComponent(event.Event, event.GameDto); com != nil {
				fmt.Fprintf(w, "event: %s\ndata: ", event.Event)
				com.Render(c.Request().Context(), c.Response().Writer)
				fmt.Fprintf(w, "\n\n")
				w.Flush()
			}
		}
	}
}

func getComponent(event string, gameDto *dto.Games) templ.Component {
	switch event {
	case "team_score":
		return gameview.GameTeamScore(*gameDto)
	case "opponent_score":
		return gameview.GameOpponentScore(*gameDto)
	}
	return nil
}
