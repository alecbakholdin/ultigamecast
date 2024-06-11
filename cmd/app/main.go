package main

import (
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"ultigamecast/internal/app/handlers"
	"ultigamecast/internal/app/middleware"
	"ultigamecast/internal/app/service"
	"ultigamecast/internal/env"
	"ultigamecast/internal/logger"
	"ultigamecast/internal/models"

	"github.com/justinas/alice"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	env.LoadEnv(".env")
	slog.SetDefault(slog.New(logger.NewHandler(&slog.HandlerOptions{})))
	db, err := sql.Open("sqlite3", "data.db")
	if err != nil {
		panic(err)
	}
	queries := models.New(db)

	authService := service.NewAuth(queries, env.MustGetenv("JWT_SECRET"))
	teamService := service.NewTeam(queries, db)
	playerService := service.NewPlayer(queries, db)
	tournamentService := service.NewTournament(queries, db)
	gameService := service.NewGame(queries, db)
	base := alice.New(
		middleware.RecoverPanic,
		middleware.LoadUser(authService),
		middleware.LoadContext(middleware.Services{Team: teamService, Tournament: tournamentService, Player: playerService, Game: gameService}),
		middleware.LogRequestAndHandleError,
	)
	if os.Getenv("USE_DELAY") != "" {
		slog.Info("Adding artificial delay to every HTTP request")
		base = base.Append(middleware.Delay)
	}
	authenticatedOnly := base.Append(middleware.GuardAuthenticated)
	adminOnly := base.Append(middleware.GuardTeamAdmin)

	// public directory
	if dir, err := filepath.Glob("./web/public/**"); err != nil {
		panic(fmt.Errorf("could not read directory web/dir: %w", err))
	} else {
		for _, f := range dir {
			http.HandleFunc(fmt.Sprintf("GET %s", strings.TrimPrefix(f, "web/public")), func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, f) })
		}
	}

	homeHandler := handlers.NewHome()
	http.Handle("GET /", base.ThenFunc(homeHandler.GetHome))

	authHandler := handlers.NewAuth(authService)
	http.Handle("GET /login", base.ThenFunc(authHandler.GetLogin))
	http.Handle("POST /login", base.ThenFunc(authHandler.PostLogin))
	http.Handle("GET /signup", base.ThenFunc(authHandler.GetSignup))
	http.Handle("POST /signup", base.ThenFunc(authHandler.PostSignup))
	http.Handle("POST /logout", base.ThenFunc(authHandler.PostLogout))

	teamHandler := handlers.NewTeam(teamService)
	http.Handle("GET /teams", authenticatedOnly.ThenFunc(teamHandler.GetTeams))
	http.Handle("GET /teams-create", authenticatedOnly.ThenFunc(teamHandler.GetTeamsCreate))
	http.Handle("POST /teams", authenticatedOnly.ThenFunc(teamHandler.PostTeams))
	http.Handle("GET /teams/{teamSlug}", base.ThenFunc(teamHandler.GetTeam))
	http.Handle("GET /teams/{teamSlug}/edit", adminOnly.ThenFunc(teamHandler.GetEdit))
	http.Handle("PUT /teams/{teamSlug}/edit", adminOnly.ThenFunc(teamHandler.PutEdit))
	http.Handle("GET /teams/{teamSlug}/edit-cancel", adminOnly.ThenFunc(teamHandler.GetCancelEdit))

	teamScheduleHandler := handlers.NewTeamSchedule(tournamentService)
	http.Handle("GET /teams/{teamSlug}/schedule", base.ThenFunc(teamScheduleHandler.Get))
	http.Handle("POST /teams/{teamSlug}/schedule", adminOnly.ThenFunc(teamScheduleHandler.Post))
	http.Handle("GET /teams/{teamSlug}/schedule-create", adminOnly.ThenFunc(teamScheduleHandler.GetModal))

	tournamentHandler := handlers.NewTournament(tournamentService)
	http.Handle("GET /teams/{teamSlug}/schedule/{tournamentSlug}", base.ThenFunc(tournamentHandler.Get))
	http.Handle("GET /teams/{teamSlug}/schedule/{tournamentSlug}/edit", adminOnly.ThenFunc(tournamentHandler.GetEdit))
	http.Handle("PUT /teams/{teamSlug}/schedule/{tournamentSlug}/edit", adminOnly.ThenFunc(tournamentHandler.PutEdit))
	http.Handle("GET /teams/{teamSlug}/schedule/{tournamentSlug}/edit-cancel", adminOnly.ThenFunc(tournamentHandler.GetEditCancel))

	log.Fatal(http.ListenAndServe("localhost:8090", nil))
}
