package setup

import (
	"time"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/core"
)

func RegisterDevParams(e *core.ServeEvent) {
	if e.App.IsDev() {
		e.Router.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				time.Sleep(time.Millisecond * 500)
				return next(c)
			}
		})
	}
}
