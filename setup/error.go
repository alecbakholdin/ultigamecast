package setup

import (
	"fmt"
	"net/http"
	"ultigamecast/view/component"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v5"
)

const (
	teamNotFoundMessage    = "Team doesn't exist"
	unexpectedErrorMessage = "Unexpected error"
)

func ErrorHandler(c echo.Context, err error) {
	code := http.StatusInternalServerError
	message := unexpectedErrorMessage
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code

		if m, ok := he.Message.(string); ok {
			message = m
		} else {
			c.Echo().Logger.Error(fmt.Errorf("unsupported message type: %s", m))
		}
	}

	c.Echo().Logger.Error(err)
	var comp templ.Component
	if c.Request().Header.Get("Hx-Request") == "true" {
		comp = component.ErrorContent(c, code, message)
	} else {
		comp = component.ErrorPage(c, code, message)
	}
	if err := comp.Render(c.Request().Context(), c.Response().Writer); err != nil {
		c.Echo().Logger.Error(err)
	}
}
