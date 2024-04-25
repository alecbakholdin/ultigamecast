package handlers

import (
	"fmt"
	"net/http"
	"ultigamecast/modelspb"
	"ultigamecast/modelspb/dto"
	"ultigamecast/repository"
	"ultigamecast/validation"
	"ultigamecast/view/component"
	view "ultigamecast/view/team/tournaments"

	"github.com/labstack/echo/v5"
)

type Tournaments struct {
	TournamentRepo *repository.Tournament
	TeamRepo       *repository.Team
}

func NewTournaments(to *repository.Tournament, te *repository.Team) *Tournaments {
	return &Tournaments{
		TournamentRepo: to,
		TeamRepo:       te,
	}
}

func (t *Tournaments) Routes(g *echo.Group) *echo.Group {
	g.GET("/tournaments", t.getTournaments)

	g.GET("/newTournament", t.getNewTournament)
	g.POST("/tournaments", t.createNewTournament)

	tournamentGroup := g.Group("/tournaments/:tournamentSlug")
	tournamentGroup.GET("/edit", t.getEditTournament)
	tournamentGroup.PUT("", t.updateTournament)
	tournamentGroup.DELETE("", t.deleteTournament)

	return tournamentGroup
}

func (t *Tournaments) getTournaments(c echo.Context) (err error) {
	var (
		teamSlug = c.PathParam("teamSlug")
	)

	if tournaments, err := t.TournamentRepo.GetAllByTeamSlug(teamSlug); err != nil {
		return echo.NewHTTPErrorWithInternal(http.StatusInternalServerError, err, "Unexpected error")
	} else {
		return view.TeamTournaments(c, teamSlug, tournaments).Render(c.Request().Context(), c.Response().Writer)
	}
}

func (t *Tournaments) getNewTournament(c echo.Context) (err error) {
	teamSlug := c.PathParam("teamSlug")
	TriggerOpenModal(c)
	return view.TournamentDialog(c, true, dto.Tournament{
		TeamSlug: teamSlug,
	}).Render(c.Request().Context(), c.Response().Writer)
}

func (t *Tournaments) createNewTournament(c echo.Context) (err error) {
	var (
		payload    dto.Tournament
		team       *modelspb.Teams
		tournament *modelspb.Tournaments
	)
	if err := dto.BindTournament(c, &payload); err != nil {
		c.Echo().Logger.Error(fmt.Errorf("error binding tournament: %s", err))
		return component.RenderToastError(c, "unexpected error")
	}
	defer renderForm(c, true, &payload)

	if exists, err := t.TournamentRepo.ExistsBySlug(payload.TeamSlug, payload.TournamentSlugNew); err != nil {
		c.Echo().Logger.Error(fmt.Errorf("error determining if tournament exists: %s", err))
		return component.RenderToastError(c, "unexpected error")
	} else if exists {
		validation.AddFieldErrorString(c, "name", "Name is already taken")
		return nil
	}

	if !validation.IsFormValid(c) {
		return nil
	}

	if team, err = t.TeamRepo.GetOneBySlug(payload.TeamSlug); err != nil {
		c.Echo().Logger.Error(fmt.Errorf("error finding team %s: %s", payload.TeamSlug, err))
		validation.AddFormErrorString(c, "unexpected error finding team")
	} else if tournament, err = t.TournamentRepo.Create(team, payload.Name, payload.TournamentSlugNew, payload.StartDt, payload.EndDt, payload.Location); err != nil {
		c.Echo().Logger.Error(fmt.Errorf("error creating tournament"))
		validation.AddFormErrorString(c, "unexpected error creating tournament")
	} else {
		TriggerCloseModal(c)
		return view.NewTournamentRow(payload.TeamSlug, tournament).Render(c.Request().Context(), c.Response().Writer)
	}
	return
}

func (t *Tournaments) getEditTournament(c echo.Context) (err error) {
	teamSlug := c.PathParam("teamSlug")
	tournamentSlug := c.PathParam("tournamentSlug")
	if teamSlug == "" || tournamentSlug == "" {
		c.Echo().Logger.Error(fmt.Errorf("teamSlug [%s] or tournamentId [%s] is empty", teamSlug, tournamentSlug))
		return component.RenderToastError(c, "unexpected error")
	}

	if tournament, err := t.TournamentRepo.GetOneBySlug(teamSlug, tournamentSlug); err != nil {
		c.Echo().Logger.Error(fmt.Errorf("error finding tournament: %s", err))
		return component.RenderToastError(c, "could not find tournament")
	} else {
		TriggerOpenModal(c)
		data := dto.DtoFromTournament(tournament)
		data.TeamSlug = teamSlug
		return view.TournamentDialog(c, false, *data).Render(c.Request().Context(), c.Response().Writer)
	}
}

func (t *Tournaments) updateTournament(c echo.Context) (err error) {
	var (
		payload    dto.Tournament
		tournament *modelspb.Tournaments
	)
	if err = dto.BindTournament(c, &payload); err != nil {
		c.Echo().Logger.Error(fmt.Errorf("error binding tournament: %s", err))
		return component.RenderToastError(c, "unexpected error")
	}
	defer renderForm(c, false, &payload)

	if payload.TournamentSlug != payload.TournamentSlugNew {
		if exists, err := t.TournamentRepo.ExistsBySlug(payload.TeamSlug, payload.TournamentSlugNew); err != nil {
			c.Echo().Logger.Error(fmt.Errorf("error determining if tournament exists: %s", err))
			validation.AddFormErrorString(c, "unexpected error")
		} else if exists {
			validation.AddFieldErrorString(c, "name", "name is already taken")
		}

	}

	if !validation.IsFormValid(c) {
		return
	}

	tournament, err = t.TournamentRepo.UpdateBySlug(payload.TeamSlug, payload.TournamentSlug, payload.Name, payload.TournamentSlugNew, payload.StartDt, payload.EndDt, payload.Location)
	if err != nil {
		validation.AddFormErrorString(c, "could not save tournament")
		return
	}

	if validation.IsFormValid(c) {
		TriggerCloseModal(c)
		return view.EditedTournamentRow(payload.TeamSlug, tournament).Render(c.Request().Context(), c.Response().Writer)
	}
	return
}

func renderForm(c echo.Context, isNew bool, payload *dto.Tournament) {
	if err := view.TournamentForm(c, isNew, *payload).Render(c.Request().Context(), c.Response().Writer); err != nil {
		c.Echo().Logger.Error(fmt.Errorf("error rendering form: %s", err))
		component.RenderToast(c, "unexpected error returning form", component.ToastSeverityError)
	}
}

func (t *Tournaments) deleteTournament(c echo.Context) (err error) {
	teamSlug := c.PathParam("teamSlug")
	tournamentSlug := c.PathParam("tournamentSlug")
	if err = t.TournamentRepo.DeleteBySlug(teamSlug, tournamentSlug); err != nil {
		c.Echo().Logger.Error(fmt.Errorf("could not delete tournament: %s", err))
		return component.RenderToastError(c, "unexpected error")
	}
	TriggerCloseModal(c)
	return
}
