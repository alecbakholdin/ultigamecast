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
	base := alice.New(middleware.LoadContext, middleware.LoadUser(authService), middleware.LogRequest)

	authHandler := handlers.NewAuth(authService)
	http.Handle("GET /login", base.ThenFunc(authHandler.GetLogin))
	http.Handle("POST /login", base.ThenFunc(authHandler.PostLogin))
	http.Handle("GET /signup", base.ThenFunc(authHandler.GetSignup))
	http.Handle("POST /signup", base.ThenFunc(authHandler.PostSignup))
	log.Fatal(http.ListenAndServe("localhost:8090", nil))
}

