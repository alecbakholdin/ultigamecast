package main

import (
	"log"
	"ultigamecast/handlers"
	"ultigamecast/repository"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

func main() {
	app := pocketbase.New()

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
        e.Router.HTTPErrorHandler = handlers.ErrorHandler

		teamRepo := repository.NewTeam(app.Dao())
		playerRepo := repository.NewPlayer(app.Dao())
		tournamentRepo := repository.NewTournament(app.Dao())

		teamHandler := handlers.NewTeam(teamRepo, playerRepo, tournamentRepo)
		teamHandler.Routes(e.Router)

		return nil
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
