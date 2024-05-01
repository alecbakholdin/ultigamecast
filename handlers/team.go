package handlers

import (
	"io"
	"net/http"
	"ultigamecast/pbmodels"
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

func (t *Team) Routes(baseGroup *echo.Group) *echo.Group {
	teamGroup := baseGroup.Group("/team/:teamsSlug")
	teamGroup.GET("", t.getTeam)
	teamGroup.GET("/logo", t.getLogo)
	return teamGroup
}

func (t *Team) getTeam(c echo.Context) (err error) {
	var (
		team     *pbmodels.Teams
		teamSlug = c.PathParam("teamsSlug")
	)

	if team, err = t.TeamRepo.FindOneBySlug(teamSlug); repository.IsNotFound(err) {
		return echo.NewHTTPErrorWithInternal(http.StatusNotFound, err, teamNotFoundMessage)
	} else if err != nil {
		return echo.NewHTTPErrorWithInternal(http.StatusInternalServerError, err, unexpectedErrorMessage)
	}
	return view.Team(c, team).Render(c.Request().Context(), c.Response().Writer)
}

func (t *Team) getLogo(c echo.Context) (err error) {
	teamSlug := c.PathParam("teamsSlug")
	if reader, err := t.TeamRepo.GetLogo(teamSlug); err != nil {
		return echo.NewHTTPErrorWithInternal(http.StatusInternalServerError, err, "error fetching logo")
	} else if _, err := io.Copy(c.Response(), reader); err != nil {
		return echo.NewHTTPErrorWithInternal(http.StatusInternalServerError, err, "error fetching logo")
	}
	return
}
