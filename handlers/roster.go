package handlers

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"ultigamecast/pbmodels"
	"ultigamecast/repository"
	"ultigamecast/service"
	"ultigamecast/validation"
	"ultigamecast/view/component"
	view "ultigamecast/view/team"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v5"
)

type Roster struct {
	PlayerService *service.Players
}

func NewRoster(p *service.Players) *Roster {
	return &Roster{
		PlayerService: p,
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
	teamSlug := c.PathParam("teamsSlug")

	if players, err := r.PlayerService.GetAllBySlug(teamSlug); err != nil {
		return echo.NewHTTPErrorWithInternal(http.StatusInternalServerError, err, "")
	} else {
		TriggerOpenModal(c)
		return view.ManageRosterDialogContent(c, players).Render(c.Request().Context(), c.Response().Writer)
	}
}

func (r *Roster) getPlayers(c echo.Context) (err error) {
	teamSlug := c.PathParam("teamsSlug")

	return view.TeamRoster(c, teamSlug).Render(c.Request().Context(), c.Response().Writer)
}

func (r *Roster) getEditPlayer(c echo.Context) (err error) {
	playerId := c.PathParam("playerId")

	if player, err := r.PlayerService.GetOneById(playerId); err != nil && !repository.IsNotFound(err) {
		return echo.NewHTTPErrorWithInternal(http.StatusInternalServerError, err, "unexpected error")
	} else {
		TriggerOpenModal(c)
		return view.EditPlayerRow(c, player).Render(c.Request().Context(), c.Response().Writer)
	}
}

func (r *Roster) getPlayerRow(c echo.Context) (err error) {
	playerId := c.PathParam("playerId")

	if player, err := r.PlayerService.GetOneById(playerId); err != nil && !repository.IsNotFound(err) {
		return echo.NewHTTPErrorWithInternal(http.StatusInternalServerError, err, "unexpected error")
	} else {
		return view.PlayerRow(c, player).Render(c.Request().Context(), c.Response().Writer)
	}
}

func (r *Roster) getRosterTable(c echo.Context) (err error) {
	var (
		teamSlug        = c.PathParam("teamsSlug")
		summaryType     = c.QueryParamDefault("type", "team")
		orderByField    = c.QueryParamDefault("orderby", "player_name")
		direction       = c.QueryParamDefault("dir", "asc")
		sortDirection   repository.SortDirection
		playerSummaries []pbmodels.PlayerTeamSummary
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
	// TODO: new summary service for stats
	return view.RosterTable(teamSlug, summaryType, playerSummaries, orderByField, sortDirection == repository.SortDirectionAsc).Render(c.Request().Context(), c.Response().Writer)
}

func (r *Roster) createPlayer(c echo.Context) (err error) {
	var (
		payload pbmodels.Players
		teamSlug = c.PathParam("teamsSlug")
	)
	if err = c.Bind(&payload); err != nil {
		c.Echo().Logger.Error(fmt.Errorf("error binding player: %s", err))
		return component.RenderToastError(c, "unexpected error")
	}
	defer renderPlayerForm(c, &payload)

	if !validation.IsFormValid(c) {
		return
	}

	if err := r.PlayerService.Create(teamSlug, &payload); err != nil {
		validation.AddFormErrorString(c, "unexpected error creating player")
	}

	if validation.IsFormValid(c) {
		payload.Order += 1
		payload.Name = ""
		return view.NewPlayerRow(c, &payload).Render(c.Request().Context(), c.Response().Writer)
	}
	return nil
}

func (r *Roster) updatePlayer(c echo.Context) (err error) {
	var (
		payload pbmodels.Players
	)
	if err := c.Bind(&payload); err != nil {
		c.Echo().Logger.Error(fmt.Errorf("error binding player: %s", err))
		return component.RenderToastError(c, "unexpected error")
	}

	defer func() {
		var component templ.Component
		if validation.IsFormValid(c) {
			component = view.PlayerRow(c, &payload)
		} else {
			component = view.EditPlayerRow(c, &payload)
		}
		if err := component.Render(c.Request().Context(), c.Response().Writer); err != nil {
			c.Echo().Logger.Error(fmt.Errorf("error rendering [validForm=%t] player row: %s", validation.IsFormValid(c), err))
		}
	}()

	if !validation.IsFormValid(c) {
		return
	}

	if err = r.PlayerService.Update(payload.Id, &payload); err != nil {
		c.Echo().Logger.Error(err)
		validation.AddFormError(c, err)
	}
	return
}

func (r *Roster) deletePlayer(c echo.Context) (err error) {
	playerId := c.PathParam("playerId")
	if err := r.PlayerService.Delete(playerId); err != nil {
		c.Echo().Logger.Error(err)
		return component.RenderToastError(c, "unexpected error")
	}
	return
}

func (r *Roster) updateRosterOrder(c echo.Context) (err error) {
	var (
		form      *multipart.Form
		playerIds []string
		ok        bool
	)
	if form, err = c.MultipartForm(); err != nil {
		c.Echo().Logger.Error(fmt.Errorf("error getting multipart form: %s", err))
		return component.RenderToastError(c, "unexpected error")
	}

	if playerIds, ok = form.Value["player_id"]; !ok {
		return
	} else if err = r.PlayerService.UpdateOrder(playerIds); err != nil {
		c.Echo().Logger.Error(err)
		validation.AddFormErrorString(c, "could not update order. Try refreshing")
	}

	return
}

func renderPlayerForm(c echo.Context, payload *pbmodels.Players) (err error) {
	err = view.CreatePlayerRow(c, payload).Render(c.Request().Context(), c.Response().Writer)
	if err != nil {
		c.Echo().Logger.Error(fmt.Errorf("error rendering CreatePlayerRow: %s", err))
	}
	return
}
