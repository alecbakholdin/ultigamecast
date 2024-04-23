package handlers

import (
	"fmt"
	"net/http"
	"ultigamecast/handlers/models"
	"ultigamecast/modelspb"
	"ultigamecast/repository"
	"ultigamecast/validation"
	"ultigamecast/view/component"
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
	g.GET("/newPlayer", r.getNewPlayer)
	g.POST("", r.createPlayer)

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

	return view.RosterTable(teamSlug, summaryType, playerSummaries, orderByField, sortDirection == repository.SortDirectionAsc).Render(c.Request().Context(), c.Response().Writer)
}

func (r *Roster) getNewPlayer(c echo.Context) (err error) {
	teamSlug := c.PathParam("teamSlug")

	return view.PlayerDialogContent(c, "New Player", view.PlayerData{TeamSlug: teamSlug}).Render(c.Request().Context(), c.Response().Writer)
}

func (r *Roster) createPlayer(c echo.Context) (err error) {
	var (
		payload models.PlayerPayload
	)
	if err = models.BindPlayer(c, &payload); err != nil {
		c.Echo().Logger.Error(fmt.Errorf("error binding player: %s", err))
		return component.RenderToastError(c, "unexpected error")
	}
	defer renderPlayerDialogContent(c, "New Player", &payload)

	if validation.IsFormValid(c) {
		MarkFormSuccess(c)
		return 
	}
	return nil
}

func renderPlayerDialogContent(c echo.Context, title string, payload *models.PlayerPayload) (err error) {
	err = view.PlayerDialogContent(c, title, payload.ToData()).Render(c.Request().Context(), c.Response().Writer)
	if err != nil {
		c.Echo().Logger.Error(fmt.Errorf("error rendering PlayerDialogContent: %s", err))
	}
	return
}
