package setup

import (
	"fmt"
	"net/http"
	"ultigamecast/view/component"
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
	if err := component.Error(c, code, message).Render(c.Request().Context(), c.Response().Writer); err != nil {
		c.Echo().Logger.Error(err)
	}
}
