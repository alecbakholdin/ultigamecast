package handlers

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"ultigamecast/handlers/models"
	"ultigamecast/modelspb"
	"ultigamecast/repository"
	"ultigamecast/validation"
	"ultigamecast/view/component"
	view "ultigamecast/view/team"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v5"
)

type Roster struct {
	PlayerRepo *repository.Player
	TeamRepo   *repository.Team
}

func NewRoster(p *repository.Player, t *repository.Team) *Roster {
	return &Roster{
		PlayerRepo: p,
		TeamRepo:   t,
	}
}

func (r *Roster) Routes(g *echo.Group) *echo.Group {
	g.GET("/manageRoster", r.getManageRoster)
	g.GET("/rosterTable", r.getRosterTable)
	g.GET("/roster", r.getPlayers)
	g.POST("/roster", r.createPlayer)
	g.PUT("/rosterOrder", r.updateRosterOrder)

	playerGroup := g.Group("/roster/:playerId")
	playerGroup.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.PathParam("playerId") == "" {
				return echo.NewHTTPError(http.StatusBadRequest, "no player specified in request")
			}
			return next(c)
		}
	})
	playerGroup.GET("/edit", r.getEditPlayer)
	playerGroup.GET("/row", r.getPlayerRow)
	playerGroup.DELETE("", r.deletePlayer)
	playerGroup.PUT("", r.updatePlayer)
	return playerGroup
}

func (r *Roster) getManageRoster(c echo.Context) (err error) {
	teamSlug := c.PathParam("teamSlug")

	if players, err := r.PlayerRepo.GetAllByTeamSlug(teamSlug); err != nil {
		return echo.NewHTTPErrorWithInternal(http.StatusInternalServerError, err, "")
	} else {
		return view.ManageRosterDialogContent(c, players).Render(c.Request().Context(), c.Response().Writer)
	}
}

func (r *Roster) getPlayers(c echo.Context) (err error) {
	teamSlug := c.PathParam("teamSlug")

	return view.TeamRoster(c, teamSlug).Render(c.Request().Context(), c.Response().Writer)
}

func (r *Roster) getEditPlayer(c echo.Context) (err error) {
	playerId := c.PathParam("playerId")

	if player, err := r.PlayerRepo.GetOneById(playerId); err != nil && !repository.IsNotFound(err) {
		return echo.NewHTTPErrorWithInternal(http.StatusInternalServerError, err, "unexpected error")
	} else {
		return view.EditPlayerRow(c, view.PlayerData{
			PlayerID:    player.Record.GetId(),
			PlayerName:  player.GetName(),
			PlayerOrder: player.GetOrder(),
		}).Render(c.Request().Context(), c.Response().Writer)
	}
}

func (r *Roster) getPlayerRow(c echo.Context) (err error) {
	playerId := c.PathParam("playerId")

	if player, err := r.PlayerRepo.GetOneById(playerId); err != nil && !repository.IsNotFound(err) {
		return echo.NewHTTPErrorWithInternal(http.StatusInternalServerError, err, "unexpected error")
	} else {
		return view.PlayerRow(c, player).Render(c.Request().Context(), c.Response().Writer)
	}
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

func (r *Roster) createPlayer(c echo.Context) (err error) {
	var (
		payload models.PlayerPayload
		player  *modelspb.Players
	)
	if err = models.BindPlayer(c, &payload); err != nil {
		c.Echo().Logger.Error(fmt.Errorf("error binding player: %s", err))
		return component.RenderToastError(c, "unexpected error")
	}
	defer renderPlayerForm(c, &payload)

	if !validation.IsFormValid(c) {
		return
	}

	if team, err := r.TeamRepo.GetOneBySlug(payload.TeamSlug); err != nil {
		c.Echo().Logger.Error(fmt.Errorf("error fetching team by slug: %s", err))
		validation.AddFormError(c, fmt.Errorf("unexpected error"))
	} else if player, err = r.PlayerRepo.Create(team, payload.Name, payload.Order); err != nil {
		c.Echo().Logger.Error(fmt.Errorf("error creating player: %s", err))
		validation.AddFormError(c, fmt.Errorf("unexpected error creating player, try refreshing"))
	}

	if validation.IsFormValid(c) {
		MarkFormSuccess(c)
		payload.Order += 1
		payload.Name = ""
		return view.NewPlayerRow(c, player).Render(c.Request().Context(), c.Response().Writer)
	}
	return nil
}

func (r *Roster) updatePlayer(c echo.Context) (err error) {
	var (
		payload models.PlayerPayload
		player  *modelspb.Players
	)
	if err := models.BindPlayer(c, &payload); err != nil {
		c.Echo().Logger.Error(fmt.Errorf("error binding player: %s", err))
		return component.RenderToastError(c, "unexpected error")
	}

	defer func(c echo.Context, payload *models.PlayerPayload, player **modelspb.Players) {
		var component templ.Component
		if validation.IsFormValid(c) {
			component = view.PlayerRow(c, *player)
		} else {
			component = view.EditPlayerRow(c, payload.ToData())
		}
		if err := component.Render(c.Request().Context(), c.Response().Writer); err != nil {
			c.Echo().Logger.Error(fmt.Errorf("error rendering [validForm=%t] player row: %s", validation.IsFormValid(c), err))
		}
	}(c, &payload, &player)

	if !validation.IsFormValid(c) {
		return
	}

	if player, err = r.PlayerRepo.Update(payload.PlayerID, payload.Name); err != nil {
		validation.AddFormError(c, err)
	}
	return
}

func (r *Roster) deletePlayer(c echo.Context) (err error) {
	playerId := c.PathParam("playerId")
	if err := r.PlayerRepo.Delete(playerId); err != nil {
		c.Echo().Logger.Error(fmt.Errorf("error deleting player %s: %s", playerId, err))
		return component.RenderToastError(c, "unexpected error")
	}
	return
}

func (r *Roster) updateRosterOrder(c echo.Context) (err error) {
	var (
		form      *multipart.Form
		playerIds []string
		players   []*modelspb.Players
		teamSlug  = c.PathParam("teamSlug")
		ok        bool
	)
	if form, err = c.MultipartForm(); err != nil {
		c.Echo().Logger.Error(fmt.Errorf("error getting multipart form: %s", err))
		return component.RenderToastError(c, "unexpected error")
	}

	if playerIds, ok = form.Value["player_id"]; !ok {
		return
	} else if players, err = r.PlayerRepo.GetAllByTeamSlug(teamSlug); err != nil {
		c.Echo().Logger.Error(fmt.Errorf("error fetching players: %s", err))
		validation.AddFormErrorString(c, "unexpected error")
	} else if len(players) != len(playerIds) {
		validation.AddFormErrorString(c, "your data is stale, please refresh and try again")
	}

	if !validation.IsFormValid(c) {
		return
	}

	if _, err = r.PlayerRepo.UpdateOrder(players, playerIds); err != nil {
		c.Echo().Logger.Error(fmt.Errorf("error updating player order: %s", err))
		validation.AddFormErrorString(c, "unexpected error")
	}

	return
}

func renderPlayerForm(c echo.Context, payload *models.PlayerPayload) (err error) {
	err = view.CreatePlayerRow(c, payload.ToData()).Render(c.Request().Context(), c.Response().Writer)
	if err != nil {
		c.Echo().Logger.Error(fmt.Errorf("error rendering CreatePlayerRow: %s", err))
	}
	return
}
