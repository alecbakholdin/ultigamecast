package main

import (
	"database/sql"
	"log"
	"log/slog"
	"net/http"
	"ultigamecast/app/env"
	"ultigamecast/app/logger"
	"ultigamecast/handlers"
	"ultigamecast/middleware"
	"ultigamecast/models"
	"ultigamecast/service"

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
	teamService := service.NewTeam(queries)
	base := alice.New(middleware.LoadContext(teamService), middleware.LoadUser(authService), middleware.LogRequest)
	mustBeAuthenticated := base.Append(middleware.GuardAuthenticated)

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
	http.Handle("GET /teams", mustBeAuthenticated.ThenFunc(teamHandler.GetTeams))
	http.Handle("POST /teams", mustBeAuthenticated.ThenFunc(teamHandler.PostTeams))

	log.Fatal(http.ListenAndServe("localhost:8090", nil))
}
