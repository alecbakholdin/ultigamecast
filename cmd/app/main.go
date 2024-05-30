package main

import (
	"database/sql"
	"log"
	"log/slog"
	"net/http"
	"os"
	"ultigamecast/internal/env"
	"ultigamecast/internal/logger"
	"ultigamecast/internal/app/middleware"
	"ultigamecast/internal/app/handlers"
	"ultigamecast/internal/app/service"
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
	base := alice.New(
		middleware.RecoverPanic,
		middleware.LoadContext(teamService),
		middleware.LoadUser(authService),
		middleware.LogRequestAndHandleError,
	)
	if os.Getenv("USE_DELAY") != "" {
		slog.Info("Adding artificial delay to every HTTP request")
		base = base.Append(middleware.Delay)
	}
	authenticatedOnly := base.Append(middleware.GuardAuthenticated)

	http.HandleFunc("GET /favicon.ico", func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, "public/favicon.ico") })
	http.HandleFunc("GET /frisbee.png", func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, "public/frisbee.png") })
	http.HandleFunc("GET /styles.css", func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, "public/styles.css") })
	http.HandleFunc("GET /theme.css", func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, "public/theme.css") })

	homeHandler := handlers.NewHome()
	http.Handle("GET /", base.ThenFunc(homeHandler.GetHome))

	authHandler := handlers.NewAuth(authService)
	http.Handle("GET /login", base.ThenFunc(authHandler.GetLogin))
	http.Handle("POST /login", base.ThenFunc(authHandler.PostLogin))
	http.Handle("GET /signup", base.ThenFunc(authHandler.GetSignup))
	http.Handle("POST /signup", base.ThenFunc(authHandler.PostSignup))
	http.Handle("POST /logout", base.ThenFunc(authHandler.PostLogout))

	teamHandler := handlers.NewTeam(teamService)
	withTeam := base.Append(middleware.LoadTeam(teamService))
	withTeamAdminOnly := withTeam.Append(middleware.GuardTeamAdmin)
	http.Handle("GET /teams", authenticatedOnly.ThenFunc(teamHandler.GetTeams))
	http.Handle("POST /teams", authenticatedOnly.ThenFunc(teamHandler.PostTeams))
	http.Handle("GET /teams/{teamSlug}", withTeam.ThenFunc(teamHandler.GetTeam))
	http.Handle("PUT /teams/{teamSlug}", withTeamAdminOnly.ThenFunc(teamHandler.PutTeam))
	http.Handle("GET /teams-edit/{teamSlug}", withTeamAdminOnly.ThenFunc(teamHandler.GetTeamsEdit))
	http.Handle("GET /teams-create", authenticatedOnly.ThenFunc(teamHandler.GetTeamsCreate))

	playerHandler := handlers.NewPlayer(playerService)
	withPlayer := base.Append(middleware.LoadPlayer(playerService))
	withPlayerAdminOnly := withPlayer.Append(middleware.GuardTeamAdmin)
	http.Handle("GET /teams/{teamSlug}/players", withTeam.ThenFunc(playerHandler.GetPlayers))
	http.Handle("POST /teams/{teamSlug}/players", withTeamAdminOnly.ThenFunc(playerHandler.PostPlayers))
	http.Handle("PUT /teams/{teamSlug}/players/{playerSlug}", withPlayerAdminOnly.ThenFunc(playerHandler.PutPlayer))
	http.Handle("POST /teams/{teamSlug}/players-order", withTeamAdminOnly.ThenFunc(playerHandler.PostPlayersOrder))

	tournamentHandler := handlers.NewTournament(tournamentService)
	withTournament := withTeam.Append(middleware.LoadTournament(tournamentService))
	withTournamentAdminOnly := withTournament.Append(middleware.GuardTeamAdmin)
	http.Handle("GET /teams/{teamSlug}/tournaments", withTeam.ThenFunc(tournamentHandler.GetTournaments))
	http.Handle("POST /teams/{teamSlug}/tournaments", withTeamAdminOnly.ThenFunc(tournamentHandler.PostTournaments))
	http.Handle("GET /teams/{teamSlug}/tournaments/{tournamentSlug}", withTournament.ThenFunc(tournamentHandler.GetTournament))
	http.Handle("POST /teams/{teamSlug}/tournaments/{tournamentSlug}/data", withTournamentAdminOnly.ThenFunc(tournamentHandler.PostData))
	http.Handle("PUT /teams/{teamSlug}/tournaments/{tournamentSlug}/data/{datumSlug}", withTournamentAdminOnly.ThenFunc(tournamentHandler.PutData))
	http.Handle("DELETE /teams/{teamSlug}/tournaments/{tournamentSlug}/data/{datumSlug}", withTournamentAdminOnly.ThenFunc(tournamentHandler.DeleteData))
	http.Handle("GET /teams/{teamSlug}/tournaments/{tournamentSlug}/row", withTournament.ThenFunc(tournamentHandler.GetTournamentRow))
	http.Handle("GET /teams/{teamSlug}/tournaments/{tournamentSlug}/date", withTournament.ThenFunc(tournamentHandler.GetDate))
	http.Handle("GET /teams/{teamSlug}/tournaments/{tournamentSlug}/edit-date", withTournamentAdminOnly.ThenFunc(tournamentHandler.GetEditDate))
	http.Handle("PUT /teams/{teamSlug}/tournaments/{tournamentSlug}/edit-date", withTournamentAdminOnly.ThenFunc(tournamentHandler.PutEditDate))

	log.Fatal(http.ListenAndServe("localhost:8090", nil))
}
