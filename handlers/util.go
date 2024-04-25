package handlers

import (
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