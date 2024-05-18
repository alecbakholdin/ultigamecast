package main

import (
	"database/sql"
	"log"
	"log/slog"
	"net/http"
	"ultigamecast/app/env"
	"ultigamecast/app/logger"
	"ultigamecast/handlers"
	"ultigamecast/models"
	"ultigamecast/service"

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
	authHandler := handlers.NewAuth(authService)
	http.HandleFunc("GET /login", authHandler.GetLogin)
	http.HandleFunc("POST /login", authHandler.PostLogin)
	http.HandleFunc("GET /signup", authHandler.GetSignup)
	http.HandleFunc("POST /signup", authHandler.PostSignup)
	log.Fatal(http.ListenAndServe("localhost:8090", nil))
}

