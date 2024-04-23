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

	tournamentGroup := g.Group("/tournaments/:tournamentId")
	tournamentGroup.GET("/edit", t.getEditTournament)
	tournamentGroup.PUT("", t.updateTournament)
	tournamentGroup.DELETE("", t.deleteTournament)

	return tournamentGroup
}

func (t *Tournaments) getTournaments(c echo.Context) (err error) {
	teamSlug := c.PathParam("teamSlug")

	if tournaments, err := t.TournamentRepo.GetAllByTeamSlug(teamSlug); err != nil {
		return echo.NewHTTPErrorWithInternal(http.StatusInternalServerError, err, "Unexpected error")
	} else {
		return view.TeamTournaments(c, teamSlug, tournaments).Render(c.Request().Context(), c.Response().Writer)
	}
}

func (t *Tournaments) getNewTournament(c echo.Context) (err error) {
	teamSlug := c.PathParam("teamSlug")
	TriggerOpenModal(c)
	return view.TournamentDialog(c, "New Tournament", teamSlug, view.TournamentData{}).Render(c.Request().Context(), c.Response().Writer)
}

func (t *Tournaments) createNewTournament(c echo.Context) (err error) {
	var (
		payload    models.TournamentPayload
		team       *modelspb.Teams
		tournament *modelspb.Tournaments
	)
	if err := models.BindTournament(c, &payload); err != nil {
		c.Echo().Logger.Error(fmt.Errorf("error binding tournament: %s", err))
		return component.RenderToastError(c, "unexpected error")
	}
	defer renderForm(c, &payload)

	if exists, err := t.TournamentRepo.ExistsBySlug(payload.TeamSlug, payload.TournamentSlug); err != nil {
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
	} else if tournament, err = t.TournamentRepo.Create(team, payload.Name, payload.TournamentSlug, payload.StartDt, payload.EndDt, payload.Location); err != nil {
		c.Echo().Logger.Error(fmt.Errorf("error creating tournament"))
		validation.AddFormErrorString(c, "unexpected error creating tournament")
	} else {
		MarkFormSuccess(c)
		return view.NewTournamentRow(payload.TeamSlug, tournament).Render(c.Request().Context(), c.Response().Writer)
	}
	return
}

func (t *Tournaments) getEditTournament(c echo.Context) (err error) {
	teamSlug := c.PathParam("teamSlug")
	tournamentId := c.PathParam("tournamentId")
	if teamSlug == "" || tournamentId == "" {
		c.Echo().Logger.Error(fmt.Errorf("teamSlug [%s] or tournamentId [%s] is empty", teamSlug, tournamentId))
		return component.RenderToastError(c, "unexpected error")
	}

	if tournament, err := t.TournamentRepo.GetOneById(tournamentId); err != nil {
		c.Echo().Logger.Error(fmt.Errorf("error finding tournament: %s", err))
		return component.RenderToastError(c, "could not find tournament")
	} else {
		TriggerOpenModal(c)
		data := view.TournamentData{
			ID:       tournament.Record.GetId(),
			Name:     tournament.GetName(),
			Start:    tournament.GetStart().Time().Format("2006-01-02"),
			End:      tournament.GetEnd().Time().Format("2006-01-02"),
			Location: tournament.GetLocation(),
		}
		return view.TournamentDialog(c, "Edit Tournament", teamSlug, data).Render(c.Request().Context(), c.Response().Writer)
	}
}

func (t *Tournaments) updateTournament(c echo.Context) (err error) {
	var (
		payload    models.TournamentPayload
		tournament *modelspb.Tournaments
	)
	if err = models.BindTournament(c, &payload); err != nil {
		c.Echo().Logger.Error(fmt.Errorf("error binding tournament: %s", err))
		return component.RenderToastError(c, "unexpected error")
	}
	defer renderForm(c, &payload)

	if payload.TournamentID == "" {
		validation.AddFormErrorString(c, "missing tournament in request")
	} else if to, err := t.TournamentRepo.GetOneBySlug(payload.TeamSlug, payload.TournamentSlug); err != nil && !repository.IsNotFound(err) {
		c.Echo().Logger.Error(fmt.Errorf("error determining if tournament exists: %s", err))
		validation.AddFormErrorString(c, "unexpected error")
	} else if to != nil && to.Record.GetId() != payload.TournamentID {
		validation.AddFieldErrorString(c, "name", "name is already taken")
	}

	if !validation.IsFormValid(c) {
		return
	}

	tournament, err = t.TournamentRepo.Update(payload.TournamentID, payload.Name, payload.TournamentSlug, payload.StartDt, payload.EndDt, payload.Location)
	if err != nil {
		validation.AddFormErrorString(c, "could not save tournament")
		return
	}

	if validation.IsFormValid(c) {
		MarkFormSuccess(c)
		return view.EditedTournamentRow(payload.TeamSlug, tournament).Render(c.Request().Context(), c.Response().Writer)
	}
	return
}

func renderForm(c echo.Context, payload *models.TournamentPayload) {
	if err := view.TournamentForm(c, payload.TeamSlug, payload.ToData()).Render(c.Request().Context(), c.Response().Writer); err != nil {
		c.Echo().Logger.Error(fmt.Errorf("error rendering form: %s", err))
		component.RenderToast(c, "unexpected error returning form", component.ToastSeverityError)
	}
}

func (t *Tournaments) deleteTournament(c echo.Context) (err error) {
	tournamentId := c.PathParam("tournamentId")
	if err = t.TournamentRepo.Delete(tournamentId); err != nil {
		c.Echo().Logger.Error(fmt.Errorf("could not delete tournament: %s", err))
		return component.RenderToastError(c, "unexpected error")
	}
	MarkFormSuccess(c)
	return
}
