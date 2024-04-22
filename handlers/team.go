package handlers

import (
	"net/http"
	"ultigamecast/modelspb"
	"ultigamecast/repository"
	view "ultigamecast/view/team"

	"github.com/labstack/echo/v5"
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
