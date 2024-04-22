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

		teamRepo := repository.NewTeam(app.Dao())
		playerRepo := repository.NewPlayer(app.Dao())
		tournamentRepo := repository.NewTournament(app.Dao())

		teamHandler := handlers.NewTeam(teamRepo, playerRepo, tournamentRepo)
		tournamentHandler := handlers.NewTournaments(tournamentRepo, teamRepo)

		baseGroup := e.Router.Group("")
		teamGroup := teamHandler.Routes(baseGroup)
		tournamentHandler.Routes(teamGroup)

		return nil
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
