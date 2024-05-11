package handlers

import (
	"fmt"
	"net/http"
	"ultigamecast/pbmodels"
	"ultigamecast/service"
	"ultigamecast/validation"
	"ultigamecast/view/component"
	view "ultigamecast/view/team/tournaments"

	"github.com/labstack/echo/v5"
)

type Tournaments struct {
	TournamentService *service.Tournaments
}

func NewTournaments(ts *service.Tournaments) *Tournaments {
	return &Tournaments{
		TournamentService: ts,
	}
}

func (t *Tournaments) Routes(g *echo.Group) *echo.Group {
	g.GET("/tournaments", t.getTournaments)

	g.GET("/newTournament", t.getNewTournament)
	g.POST("/tournaments", t.createNewTournament)

	tournamentGroup := g.Group("/tournaments/:tournamentsSlug")
	tournamentGroup.GET("/edit", t.getEditTournament)
	tournamentGroup.PUT("", t.updateTournament)
	tournamentGroup.DELETE("", t.deleteTournament)

	return tournamentGroup
}

func (t *Tournaments) getTournaments(c echo.Context) (err error) {
	teamSlug := c.PathParam("teamsSlug")
	if tg, err := t.TournamentService.GetTournamentsWithGamesByTeamSlug(teamSlug); err != nil {
		return echo.NewHTTPErrorWithInternal(http.StatusInternalServerError, err, "unexpected error")
	} else {
		return view.TeamTournaments(c, teamSlug, tg).Render(c.Request().Context(), c.Response().Writer)
	}
}

func (t *Tournaments) getNewTournament(c echo.Context) (err error) {
	TriggerOpenModal(c)
	return view.TournamentDialog(c, true, &pbmodels.Tournaments{}).Render(c.Request().Context(), c.Response().Writer)
}

func (t *Tournaments) createNewTournament(c echo.Context) (err error) {
	var (
		tournament     *pbmodels.Tournaments = &pbmodels.Tournaments{}
		teamSlug                             = c.PathParam("teamsSlug")
	)
	if err := c.Bind(tournament); err != nil {
		c.Echo().Logger.Error(fmt.Errorf("error binding tournament: %s", err))
		return component.RenderToastError(c, "unexpected error")
	}
	defer renderForm(c, true, tournament)

	tournament.Slug = ConvertToSlug(tournament.Name)
	t.TournamentService.ValidateBasicFields(c, tournament)
	t.TournamentService.ValidateSlugChange(c, teamSlug, "", tournament.Slug)

	if !validation.IsFormValid(c) {
		return nil
	}

	if t.TournamentService.Create(teamSlug, tournament); err != nil {
		c.Echo().Logger.Error(err)
		validation.AddFormErrorString(c, "unexpected error creating tournament")
		return
	}
	TriggerCloseModal(c)
	return view.NewTournamentRow(c, tournament).Render(c.Request().Context(), c.Response().Writer)
}

func (t *Tournaments) getEditTournament(c echo.Context) (err error) {
	teamSlug := c.PathParam("teamsSlug")
	tournamentSlug := c.PathParam("tournamentsSlug")
	if teamSlug == "" || tournamentSlug == "" {
		c.Echo().Logger.Error(fmt.Errorf("teamSlug [%s] or tournamentId [%s] is empty", teamSlug, tournamentSlug))
		return component.RenderToastError(c, "unexpected error")
	}

	if tournament, err := t.TournamentService.GetOneBySlug(teamSlug, tournamentSlug); err != nil {
		c.Echo().Logger.Error(fmt.Errorf("error finding tournament: %s", err))
		return component.RenderToastError(c, "could not find tournament")
	} else {
		TriggerOpenModal(c)
		return view.TournamentDialog(c, false, tournament).Render(c.Request().Context(), c.Response().Writer)
	}
}

func (t *Tournaments) updateTournament(c echo.Context) (err error) {
	var (
		tournament        pbmodels.Tournaments
		teamSlug          = c.PathParam("teamsSlug")
		oldTournamentSlug = c.PathParam("tournamentsSlug")
	)
	if err = c.Bind(&tournament); err != nil {
		c.Echo().Logger.Error(fmt.Errorf("error binding tournament: %s", err))
		return component.RenderToastError(c, "unexpected error")
	}
	defer renderForm(c, false, &tournament)

	tournament.Slug = ConvertToSlug(tournament.Name)
	t.TournamentService.ValidateBasicFields(c, &tournament)
	t.TournamentService.ValidateSlugChange(c, teamSlug, oldTournamentSlug, tournament.Slug)

	if !validation.IsFormValid(c) {
		return
	}

	if err = t.TournamentService.UpdateBySlug(teamSlug, oldTournamentSlug, &tournament); err != nil {
		validation.AddFormErrorString(c, "could not save tournament")
		return
	}

	TriggerCloseModal(c)
	return view.EditedTournamentRow(c, &tournament).Render(c.Request().Context(), c.Response().Writer)
}

func renderForm(c echo.Context, isNew bool, payload *pbmodels.Tournaments) {
	if err := view.TournamentForm(c, isNew, payload).Render(c.Request().Context(), c.Response().Writer); err != nil {
		c.Echo().Logger.Error(fmt.Errorf("error rendering form: %s", err))
		component.RenderToast(c, "unexpected error returning form", component.ToastSeverityError)
	}
}

func (t *Tournaments) deleteTournament(c echo.Context) (err error) {
	teamSlug := c.PathParam("teamsSlug")
	tournamentSlug := c.PathParam("tournamentsSlug")
	if err = t.TournamentService.Delete(teamSlug, tournamentSlug); err != nil {
		c.Echo().Logger.Error(err)
		return component.RenderToastError(c, "unexpected error")
	}
	TriggerCloseModal(c)
	return
}
