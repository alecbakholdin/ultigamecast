package handlers

import (
	"fmt"
	"net/http"
	"ultigamecast/pbmodels"
	"ultigamecast/repository"
	"ultigamecast/service"
	"ultigamecast/validation"
	"ultigamecast/view/component"
	view "ultigamecast/view/team/tournaments"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/tools/types"
)

type Tournaments struct {
	TournamentService *service.Tournaments
	TournamentRepo    *repository.Tournament
	GameRepo          *repository.Game
	TeamRepo          *repository.Team
}

func NewTournaments(ts *service.Tournaments, to *repository.Tournament, g *repository.Game, te *repository.Team) *Tournaments {
	return &Tournaments{
		TournamentService: ts,
		TournamentRepo:    to,
		GameRepo:          g,
		TeamRepo:          te,
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
		team           *pbmodels.Teams
		tournament     *pbmodels.Tournaments = &pbmodels.Tournaments{}
		teamSlug                             = c.PathParam("teamsSlug")
		tournamentSlug string
	)
	if err := c.Bind(tournament); err != nil {
		c.Echo().Logger.Error(fmt.Errorf("error binding tournament: %s", err))
		return component.RenderToastError(c, "unexpected error")
	}
	defer renderForm(c, true, tournament)

	tournamentSlug = ConvertToSlug(tournament.Name)
	if exists, err := t.TournamentRepo.ExistsBySlug(teamSlug, tournamentSlug); err != nil {
		c.Echo().Logger.Error(fmt.Errorf("error determining if tournament exists: %s", err))
		return component.RenderToastError(c, "unexpected error")
	} else if exists {
		validation.AddFieldErrorString(c, "name", "Name is already taken")
		return nil
	}

	if tournament.Start, err = tournament.GetStartDt(); err != nil {
		validation.AddFieldErrorString(c, "start", "invalid start date format")
		c.Echo().Logger.Error(fmt.Errorf("error parsing start: %s", err))
	}
	if tournament.End, err = tournament.GetEndDt(); err != nil {
		validation.AddFieldErrorString(c, "end", "invalid end date format")
		c.Echo().Logger.Error(fmt.Errorf("error parsing end: %s", err))
	}

	if !validation.IsFormValid(c) {
		return nil
	}

	if team, err = t.TeamRepo.FindOneBySlug(teamSlug); err != nil {
		c.Echo().Logger.Error(fmt.Errorf("error finding team %s: %s", teamSlug, err))
		validation.AddFormErrorString(c, "unexpected error finding team")
	} else if tournament, err = t.TournamentRepo.Create(team.Id, tournament.Name, tournamentSlug, tournament.Start, tournament.End, tournament.Location); err != nil {
		c.Echo().Logger.Error(fmt.Errorf("error creating tournament"))
		validation.AddFormErrorString(c, "unexpected error creating tournament")
	} else {
		TriggerCloseModal(c)
		return view.NewTournamentRow(tournament.Slug, tournament).Render(c.Request().Context(), c.Response().Writer)
	}
	return
}

func (t *Tournaments) getEditTournament(c echo.Context) (err error) {
	teamSlug := c.PathParam("teamsSlug")
	tournamentSlug := c.PathParam("tournamentsSlug")
	if teamSlug == "" || tournamentSlug == "" {
		c.Echo().Logger.Error(fmt.Errorf("teamSlug [%s] or tournamentId [%s] is empty", teamSlug, tournamentSlug))
		return component.RenderToastError(c, "unexpected error")
	}

	if tournament, err := t.TournamentRepo.GetOneBySlug(teamSlug, tournamentSlug); err != nil {
		c.Echo().Logger.Error(fmt.Errorf("error finding tournament: %s", err))
		return component.RenderToastError(c, "could not find tournament")
	} else {
		TriggerOpenModal(c)
		return view.TournamentDialog(c, false, tournament).Render(c.Request().Context(), c.Response().Writer)
	}
}

func (t *Tournaments) updateTournament(c echo.Context) (err error) {
	var (
		tournament *pbmodels.Tournaments = &pbmodels.Tournaments{}
		teamSlug                         = c.PathParam("teamsSlug")
		startDt    types.DateTime
		endDt      types.DateTime
	)
	if err = c.Bind(tournament); err != nil {
		c.Echo().Logger.Error(fmt.Errorf("error binding tournament: %s", err))
		return component.RenderToastError(c, "unexpected error")
	}
	defer renderForm(c, false, tournament)

	newSlug := ConvertToSlug(tournament.Name)
	if tournament.Slug != ConvertToSlug(tournament.Name) {
		if exists, err := t.TournamentRepo.ExistsBySlug(teamSlug, newSlug); err != nil {
			c.Echo().Logger.Error(fmt.Errorf("error determining if tournament exists: %s", err))
			validation.AddFormErrorString(c, "unexpected error")
		} else if exists {
			validation.AddFieldErrorString(c, "name", "name is already taken")
		}
	}

	if !validation.IsFormValid(c) {
		return
	}
	tournament, err = t.TournamentRepo.UpdateBySlug(teamSlug, tournament.Slug, tournament.Name, newSlug, startDt, endDt, tournament.Location)
	if err != nil {
		validation.AddFormErrorString(c, "could not save tournament")
		return
	}

	if validation.IsFormValid(c) {
		TriggerCloseModal(c)
		return view.EditedTournamentRow(teamSlug, tournament).Render(c.Request().Context(), c.Response().Writer)
	}
	return
}

func validateTournament(c echo.Context, t *pbmodels.Tournaments) {

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
	if err = t.TournamentRepo.DeleteBySlug(teamSlug, tournamentSlug); err != nil {
		c.Echo().Logger.Error(fmt.Errorf("could not delete tournament: %s", err))
		return component.RenderToastError(c, "unexpected error")
	}
	TriggerCloseModal(c)
	return
}
