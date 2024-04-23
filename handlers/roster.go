package handlers

import (
	"fmt"
	"net/http"
	"ultigamecast/modelspb"
	"ultigamecast/repository"
	view "ultigamecast/view/team"

	"github.com/labstack/echo/v5"
)

type Roster struct {
	PlayerRepo *repository.Player
}

func NewRoster(p *repository.Player) *Roster {
	return &Roster{
		PlayerRepo: p,
	}
}

func (r *Roster) Routes(g *echo.Group) *echo.Group {
	g.GET("/roster", r.getPlayers)
	g.GET("/rosterTable", r.getRosterTable)

	playerGroup := g.Group("/roster/:playerId")
	return playerGroup
}

func (r *Roster) getPlayers(c echo.Context) (err error) {
	teamSlug := c.PathParam("teamSlug")

	return view.TeamRoster(c, teamSlug).Render(c.Request().Context(), c.Response().Writer)
}

func (r *Roster) getRosterTable(c echo.Context) (err error) {
	var (
		teamSlug        = c.PathParam("teamSlug")
		summaryType     = c.QueryParamDefault("type", "team")
		orderByField    = c.QueryParamDefault("orderby", "player_name")
		direction       = c.QueryParamDefault("dir", "asc")
		sortDirection   repository.SortDirection
		playerSummaries []modelspb.PlayerSummary
	)

	if direction != "asc" && direction != "desc" {
		return echo.NewHTTPError(400, "Invalid sort direction "+direction)
	} else if summaryType != "team" && summaryType != "tournament" && summaryType != "game" {
		return echo.NewHTTPError(400, "Invalid summary type"+summaryType)
	}

	if direction == "asc" {
		sortDirection = repository.SortDirectionAsc
	} else {
		sortDirection = repository.SortDirectionDesc
	}

	switch summaryType {
	case "team":
		playerSummaries, err = r.PlayerRepo.GetPlayerTeamSummariesByTeamSlug(teamSlug, orderByField, sortDirection)
	default:
		return echo.NewHTTPError(http.StatusNotImplemented, fmt.Sprintf("Roster table for %s is not implemented yet", summaryType))
	}

	if err != nil {
		return echo.NewHTTPErrorWithInternal(http.StatusInternalServerError, err, "unexpected error")
	}

	return view.RosterTable(teamSlug, playerSummaries, orderByField, sortDirection == repository.SortDirectionAsc).Render(c.Request().Context(), c.Response().Writer)
}
