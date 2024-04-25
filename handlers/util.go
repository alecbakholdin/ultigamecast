package handlers

import (
	"regexp"
	"strings"

	"github.com/labstack/echo/v5"
)

type Handler interface {
	Routes(g *echo.Group) *echo.Group
}

func TriggerCloseModal(c echo.Context) {
	c.Response().Header().Set("HX-Trigger", "closemodal")
}

func TriggerOpenModal(c echo.Context) {
	c.Response().Header().Set("HX-Trigger", "openmodal")
}

var whitespaceRegex = regexp.MustCompile(`[\s+]`)

func ConvertToSlug(s string) string {
	noWhitespace := whitespaceRegex.ReplaceAllString(s, "-")
	return strings.ToLower(noWhitespace)
}
