package handlers

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/labstack/echo/v5"
)

type Handler interface {
	Routes(g *echo.Group) *echo.Group
}

func MarkFormSuccess(c echo.Context) {
	c.Response().Header().Set("HX-Trigger", "formsuccess")
	c.Response().WriteHeader(http.StatusOK)
}

var whitespaceRegex = regexp.MustCompile(`[\s+]`)

func ConvertToSlug(s string) string {
	noWhitespace := whitespaceRegex.ReplaceAllString(s, "-")
	return strings.ToLower(noWhitespace)
}
