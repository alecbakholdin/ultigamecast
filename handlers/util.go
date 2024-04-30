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
	c.Response().Header().Add("HX-Trigger", "closemodal")
}

func TriggerOpenModal(c echo.Context) {
	c.Response().Header().Add("HX-Trigger", "openmodal")
}

var whitespaceRegex = regexp.MustCompile(`\s+`)
var nonAlphaRegex = regexp.MustCompile(`[^\w-_]`)

func ConvertToSlug(s string) string {
	noWhitespace := whitespaceRegex.ReplaceAllString(s, "-")
	noNonAlpha := nonAlphaRegex.ReplaceAllString(noWhitespace, "")
	return strings.ToLower(noNonAlpha)
}
