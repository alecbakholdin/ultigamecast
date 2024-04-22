package handlers

import (
	"cmp"
	"fmt"
	"net/http"
	"ultigamecast/modelspb"
	"ultigamecast/repository"
	"ultigamecast/validation"
	"ultigamecast/view/component"
	view "ultigamecast/view/team"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/tools/types"
)

const (
	teamNotFoundMessage    = "Team doesn't exist"
	unexpectedErrorMessage = "Unexpected error"
)

type Team struct {
	PlayerRepo     *repository.Player
	TournamentRepo *repository.Tournament
	TeamRepo       *repository.Team
}

func NewTeam(t *repository.Team, p *repository.Player, to *repository.Tournament) *Team {
	return &Team{
		TeamRepo:       t,
		PlayerRepo:     p,
		TournamentRepo: to,
	}
}

func (t *Team) Routes(g *echo.Group) (*echo.Group) {
	group := g.Group("/team/:teamSlug")
	group.GET("", t.getTeam)

	return group
}

type TournamentPayload struct {
	TeamSlug       string `path:"teamSlug"`
	Team           modelspb.Teams
	TournamentID   string `form:"tournament_id"`
	Name           string `form:"name"`
	TournamentSlug string `path:"tournamentSlug"`
	Start          string `form:"start"`
	StartDt        types.DateTime
	End            string `form:"end"`
	EndDt          types.DateTime
	Location       string `form:"location"`
}

func (t *Team) getTeam(c echo.Context) (err error) {
	var (
		team     *modelspb.Teams
		teamSlug = c.PathParam("teamSlug")
	)

	if team, err = t.TeamRepo.GetOneBySlug(teamSlug); repository.IsNotFound(err) {
		return echo.NewHTTPErrorWithInternal(http.StatusNotFound, err, teamNotFoundMessage)
	} else if err != nil {
		return echo.NewHTTPErrorWithInternal(http.StatusInternalServerError, err, unexpectedErrorMessage)
	}
	return view.Team(c, team).Render(c.Request().Context(), c.Response().Writer)
}

func (t *Team) getTournaments(c echo.Context) (err error) {
	var (
		tournaments []*modelspb.Tournaments
		teamSlug    = c.PathParam("teamSlug")
	)

	if tournaments, err = t.TournamentRepo.GetAllByTeamSlug(teamSlug); err != nil && !repository.IsNotFound(err) {
		return echo.NewHTTPErrorWithInternal(http.StatusInternalServerError, err, unexpectedErrorMessage)
	}
	return view.TeamTournaments(c, teamSlug, tournaments).Render(c.Request().Context(), c.Response().Writer)
}

func (t *Team) postTournaments(c echo.Context) (err error) {
	var (
		team       *modelspb.Teams
		tournament *modelspb.Tournaments
		payload    *TournamentPayload
	)

	if payload, err = t.validateAndBindTournament(c); err != nil {
		c.Echo().Logger.Error(err)
		return component.RenderToast(c, err.Error(), component.ToastSeverityError)
	}

	// check for duplicates
	if t, err := t.TournamentRepo.GetOneBySlug(payload.TeamSlug, payload.TournamentSlug); err != nil && !repository.IsNotFound(err) {
		c.Echo().Logger.Error(err)
		validation.AddFormErrorString(c, "unexpected error occurred creating tournament")
	} else if t != nil {
		validation.AddFieldErrorString(c, "name", "a tournament with this name already exists")
	}

	if team, err = t.TeamRepo.GetOneBySlug(payload.TeamSlug); repository.IsNotFound(err) {
		return component.RenderToast(c, teamNotFoundMessage, component.ToastSeverityError)
	} else if err != nil {
		return component.RenderToast(c, unexpectedErrorMessage, component.ToastSeverityError)
	}

	if !validation.IsFormValid(c) {
		return view.CreateTournamentForm(c, payload.TeamSlug).Render(c.Request().Context(), c.Response().Writer)
	}

	// create and return values
	if tournament, err = t.TournamentRepo.Create(team, payload.Name, payload.TournamentSlug, payload.StartDt, payload.EndDt, payload.Location); err != nil {
		c.Echo().Logger.Error(err)
		validation.AddFormErrorString(c, "unexpected error occurred creating team")
		return view.CreateTournamentForm(c, payload.TeamSlug).Render(c.Request().Context(), c.Response().Writer)
	}

	return cmp.Or(
		MarkFormSuccess(c),
		view.CreateTournamentForm(c, payload.TeamSlug).Render(c.Request().Context(), c.Response().Writer),
		view.NewTournamentRow(payload.TeamSlug, tournament).Render(c.Request().Context(), c.Response().Writer),
	)
}

func (t *Team) updateTournament(c echo.Context) (err error) {
	var (
		team       *modelspb.Teams
		tournament *modelspb.Tournaments
		payload    *TournamentPayload
	)

	if payload, err = t.validateAndBindTournament(c); err != nil {
		c.Echo().Logger.Error(err)
		return component.RenderToast(c, err.Error(), component.ToastSeverityError)
	}

	if team, err = t.TeamRepo.GetOneBySlug(payload.TeamSlug); repository.IsNotFound(err) {
		return component.RenderToast(c, teamNotFoundMessage, component.ToastSeverityError)
	} else if err != nil {
		return component.RenderToast(c, unexpectedErrorMessage, component.ToastSeverityError)
	}

	if !validation.IsFormValid(c) {
		return view.CreateTournamentForm(c, payload.TeamSlug).Render(c.Request().Context(), c.Response().Writer)
	}

	// create and return values
	if tournament, err = t.TournamentRepo.Create(team, payload.Name, payload.TournamentSlug, payload.StartDt, payload.EndDt, payload.Location); err != nil {
		c.Echo().Logger.Error(err)
		validation.AddFormErrorString(c, "unexpected error occurred creating team")
		return view.CreateTournamentForm(c, payload.TeamSlug).Render(c.Request().Context(), c.Response().Writer)
	}

	return cmp.Or(
		MarkFormSuccess(c),
		view.CreateTournamentForm(c, payload.TeamSlug).Render(c.Request().Context(), c.Response().Writer),
		view.NewTournamentRow(payload.TeamSlug, tournament).Render(c.Request().Context(), c.Response().Writer),
	)
}

func (t *Team) validateAndBindTournament(c echo.Context) (payload *TournamentPayload, err error) {
	payload = new(TournamentPayload)
	if err := c.Bind(payload); err != nil {
		return nil, fmt.Errorf("data was improperly formatted")
	}

	payload.TournamentSlug = ConvertToSlug(payload.Name)
	if payload.Name == "" {
		validation.AddFieldErrorString(c, "name", "name cannot be empty")
	} 
	if payload.StartDt, err = types.ParseDateTime(payload.Start); payload.Start != "" && err != nil {
		validation.AddFieldErrorString(c, "start", "invalid date format")
	}
	if payload.EndDt, err = types.ParseDateTime(payload.End); payload.End != "" && err != nil {
		validation.AddFieldErrorString(c, "end", "invalid date format")
	}
	if payload.StartDt.Time() != payload.EndDt.Time() && !payload.EndDt.Time().After(payload.StartDt.Time()) {
		validation.AddFieldErrorString(c, "end", "end date must be after start date")
	}

	return payload, nil
}

func (t *Team) getRoster(c echo.Context) (err error) {
	var (
		team     *modelspb.Teams
		players  []*modelspb.Players
		teamSlug = c.PathParam("teamSlug")
	)

	if team, err = t.TeamRepo.GetOneBySlug(teamSlug); repository.IsNotFound(err) {
		return echo.NewHTTPErrorWithInternal(http.StatusNotFound, err, teamNotFoundMessage)
	} else if err != nil {
		return echo.NewHTTPErrorWithInternal(http.StatusInternalServerError, err, unexpectedErrorMessage)
	}

	if players, err = t.PlayerRepo.GetAllByTeamSlug(teamSlug); err != nil && !repository.IsNotFound(err) {
		return echo.NewHTTPErrorWithInternal(http.StatusInternalServerError, err, unexpectedErrorMessage)
	}
	return view.TeamRoster(c, team, players).Render(c.Request().Context(), c.Response().Writer)
}
