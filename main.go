package main

import (
	"log"
	"ultigamecast/modelspb"
	"ultigamecast/view/team"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

func main() {
    app := pocketbase.New()

    app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
        e.Router.GET("/team", func(c echo.Context) error {
            return team.NoTeam().Render(c.Request().Context(), c.Response().Writer)
        })
        e.Router.GET("/team/:teamSlug", func(c echo.Context) error {
            teamRecord, err := app.Dao().FindFirstRecordByData("teams", "slug", c.PathParam("teamSlug"))
            if err != nil {
                panic(err)
            }
            return team.Team(modelspb.Teams{Record: teamRecord}).Render(c.Request().Context(), c.Response().Writer)
        }, apis.ActivityLogger(app))
        
        return nil
    })

    if err := app.Start(); err != nil {
        log.Fatal(err)
    }
}