package models

import (
	"ultigamecast/validation"
	view "ultigamecast/view/team"

	"github.com/labstack/echo/v5"
)

type PlayerPayload struct {
	TeamSlug string `param:"teamSlug"`
	PlayerID string `param:"playerId"`
	Name     string `form:"name"`
	Order    int    `form:"order"`
}

func BindPlayer(c echo.Context, payload *PlayerPayload) error {
	if err := c.Bind(payload); err != nil {
		return err
	}

	if payload.Name == "" {
		validation.AddFieldErrorString(c, "name", "name cannot be empty")
	}

	return nil
}

func (p *PlayerPayload) ToData() view.PlayerData {
	return view.PlayerData{
		TeamSlug:    p.TeamSlug,
		PlayerID:    p.PlayerID,
		PlayerName:  p.Name,
		PlayerOrder: p.Order,
	}
}
