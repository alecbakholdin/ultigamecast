package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"ultigamecast/handlers"
	"ultigamecast/repository"
	"ultigamecast/service"
	"ultigamecast/setup"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/spf13/cobra"
)

func main() {
	app := pocketbase.New()

	BindAppHooks(app)

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}

func BindAppHooks(app *pocketbase.PocketBase) {
	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.HTTPErrorHandler = setup.ErrorHandler
		e.Router.GET("/*", apis.StaticDirectoryHandler(os.DirFS("./public"), false))
		logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))
		e.Router.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
			LogStatus:   true,
			LogURI:      true,
			LogError:    true,
			HandleError: true, // forwards error to the global error handler, so it can decide appropriate status code
			LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
				if v.Error == nil {
					logger.LogAttrs(context.Background(), slog.LevelInfo, "REQUEST",
						slog.String("uri", v.URI),
						slog.Int("status", v.Status),
					)
				} else {
					logger.LogAttrs(context.Background(), slog.LevelError, "REQUEST_ERROR",
						slog.String("uri", v.URI),
						slog.Int("status", v.Status),
						slog.String("err", v.Error.Error()),
					)
				}
				return nil
			},
		}))

		setup.RegisterDevParams(e)

		transactionRepo := repository.NewTransaction(app)
		teamRepo := repository.NewTeam(app)
		playerRepo := repository.NewPlayer(app)
		tournamentRepo := repository.NewTournament(app)
		gameRepo := repository.NewGame(app)
		eventsRepo := repository.NewEvents(app)

		tournamentsService := service.NewTournaments(tournamentRepo, gameRepo, teamRepo)
		gamesService := service.NewGames(tournamentRepo, gameRepo)
		playersService := service.NewPlayers(playerRepo, teamRepo)
		eventsService := service.NewEvents(eventsRepo)
		liveGamesService := service.NewLiveGames(gameRepo, eventsRepo, transactionRepo)

		teamHandler := handlers.NewTeam(teamRepo, playerRepo, tournamentRepo)
		tournamentHandler := handlers.NewTournaments(tournamentsService)
		rosterHandler := handlers.NewRoster(playersService)
		gameHandler := handlers.NewGames(teamRepo, tournamentRepo, gamesService, eventsService, liveGamesService)

		baseGroup := e.Router.Group("")
		teamGroup := teamHandler.Routes(baseGroup)
		rosterHandler.Routes(teamGroup)
		tournamentGroup := tournamentHandler.Routes(teamGroup)
		gameHandler.Routes(tournamentGroup)

		return nil
	})

	app.RootCmd.AddCommand(&cobra.Command{
		Use: "types",
		Run: func(cmd *cobra.Command, args []string) {
			setup.CreateTypes(app)
		},
	})
}
