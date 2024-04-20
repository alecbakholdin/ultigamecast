package handlers

import (
	"cmp"
	"net/http"
	"strings"
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

func (t *Team) Routes(e *echo.Echo) {
	group := e.Group("/team")
	group.GET("/:teamSlug", t.getTeam)
	group.GET("/:teamSlug/tournaments", t.getTournaments)
	group.POST("/:teamSlug/tournaments", t.postTournaments)
	group.GET("/:teamSlug/roster", t.getRoster)
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
		team        *modelspb.Teams
		tournaments []*modelspb.Tournaments
		teamSlug    = c.PathParam("teamSlug")
	)

	if team, err = t.TeamRepo.GetOneBySlug(teamSlug); repository.IsNotFound(err) {
		return echo.NewHTTPErrorWithInternal(http.StatusNotFound, err, teamNotFoundMessage)
	} else if err != nil {
		return echo.NewHTTPErrorWithInternal(http.StatusInternalServerError, err, unexpectedErrorMessage)
	}

	if tournaments, err = t.TournamentRepo.GetAllByTeamSlug(teamSlug); err != nil && !repository.IsNotFound(err) {
		return echo.NewHTTPErrorWithInternal(http.StatusInternalServerError, err, unexpectedErrorMessage)
	}
	return view.TeamTournaments(c, team, tournaments).Render(c.Request().Context(), c.Response().Writer)
}

func (t *Team) postTournaments(c echo.Context) (err error) {
	var (
		team       *modelspb.Teams
		tournament *modelspb.Tournaments
		teamSlug   = c.PathParam("teamSlug")

		name           = strings.TrimSpace(c.FormValue("name"))
		tournamentSlug = ConvertToSlug(name)

		start   = c.FormValue("start")
		startDt types.DateTime

		end   = c.FormValue("end")
		endDt types.DateTime

		location = strings.TrimSpace(c.FormValue("location"))
	)

	if team, err = t.TeamRepo.GetOneBySlug(teamSlug); repository.IsNotFound(err) {
		return component.RenderToast(c, teamNotFoundMessage, component.ToastSeverityError)
	} else if err != nil {
		return component.RenderToast(c, unexpectedErrorMessage, component.ToastSeverityError)
	}

	// field value validations
	if name == "" {
		validation.AddFieldErrorString(c, "name", "name cannot be empty")
	}
	if startDt, err = types.ParseDateTime(start); start != "" && err != nil {
		validation.AddFieldErrorString(c, "start", "invalid date format")
	}
	if endDt, err = types.ParseDateTime(end); end != "" && err != nil {
		validation.AddFieldErrorString(c, "end", "invalid date format")
	}
	if startDt.Time() != endDt.Time() && !endDt.Time().After(startDt.Time()) {
		validation.AddFieldErrorString(c, "end", "end date must be after start date")
	}

	if !validation.IsFormValid(c) {
		return view.CreateTournamentForm(c, team).Render(c.Request().Context(), c.Response().Writer)
	}

	// determine if duplicate
	if t, err := t.TournamentRepo.GetOneBySlug(teamSlug, tournamentSlug); err != nil && !repository.IsNotFound(err) {
		c.Echo().Logger.Error(err)
		validation.AddFormErrorString(c, "unexpected error occurred creating team")
	} else if t != nil {
		validation.AddFieldErrorString(c, "name", "a tournament with this name already exists")
	}

	if !validation.IsFormValid(c) {
		return view.CreateTournamentForm(c, team).Render(c.Request().Context(), c.Response().Writer)
	}

	// create and return values
	if tournament, err = t.TournamentRepo.Create(team, name, tournamentSlug, startDt, endDt, location); err != nil {
		c.Echo().Logger.Error(err)
		validation.AddFormErrorString(c, "unexpected error occurred creating team")
		return view.CreateTournamentForm(c, team).Render(c.Request().Context(), c.Response().Writer)
	}

	return cmp.Or(
		MarkFormSuccess(c),
		view.CreateTournamentForm(c, team).Render(c.Request().Context(), c.Response().Writer),
		view.NewTournamentRow(team, tournament).Render(c.Request().Context(), c.Response().Writer),
	)
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
