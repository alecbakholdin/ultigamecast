package main

import (
	"log"
	"os"
	"ultigamecast/handlers"
	"ultigamecast/repository"
	"ultigamecast/setup"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

func main() {
	app := pocketbase.New()

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.HTTPErrorHandler = setup.ErrorHandler
		e.Router.GET("/*", apis.StaticDirectoryHandler(os.DirFS("./public"), false))

		setup.RegisterDevParams(e)

		teamRepo := repository.NewTeam(app)
		playerRepo := repository.NewPlayer(app)
		tournamentRepo := repository.NewTournament(app)
		liveGameRepo := repository.NewLiveGame()
		gameRepo := repository.NewGame(app, liveGameRepo)

		teamHandler := handlers.NewTeam(teamRepo, playerRepo, tournamentRepo)
		tournamentHandler := handlers.NewTournaments(tournamentRepo, teamRepo)
		rosterHandler := handlers.NewRoster(playerRepo, teamRepo)
		gameHandler := handlers.NewGames(teamRepo, tournamentRepo, gameRepo)
		gameDetailsHandler := handlers.NewGameDetails(gameRepo, liveGameRepo)

		baseGroup := e.Router.Group("")
		teamGroup := teamHandler.Routes(baseGroup)
		rosterHandler.Routes(teamGroup)
		tournamentGroup := tournamentHandler.Routes(teamGroup)
		gameGroup := gameHandler.Routes(tournamentGroup)
		gameDetailsHandler.Routes(gameGroup)

		return nil
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
