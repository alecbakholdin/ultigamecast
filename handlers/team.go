package handlers

import (
	"fmt"
	"log"
	"ultigamecast/modelspb"
	"ultigamecast/repository"
	view "ultigamecast/view/team"

	"github.com/labstack/echo/v5"
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
	group.GET("/:teamSlug/roster", t.getRoster)
}

func (t *Team) getTeam(c echo.Context) (err error) {
	var (
		team     *modelspb.Teams
		teamSlug = c.PathParam("teamSlug")
	)

	if team, err = t.TeamRepo.GetOneBySlug(teamSlug); err != nil {
		return err
	}
	return view.Team(c, team).Render(c.Request().Context(), c.Response().Writer)
}

func (t *Team) getTournaments(c echo.Context) (err error) {
	var (
		team        *modelspb.Teams
		tournaments []*modelspb.Tournaments
		teamSlug    = c.PathParam("teamSlug")
	)
	log.Println(teamSlug)

	if team, err = t.TeamRepo.GetOneBySlug(teamSlug); err != nil {
		return err
	}
	if tournaments, err = t.TournamentRepo.GetAllByTeamSlug(teamSlug); err != nil {
		return fmt.Errorf("error fetching players for %s: %s", teamSlug, err)
	}
	log.Println(team, tournaments)
	return view.TeamTournaments(c, team, tournaments).Render(c.Request().Context(), c.Response().Writer)
}

func (t *Team) getRoster(c echo.Context) (err error) {
	var (
		team     *modelspb.Teams
		players  []*modelspb.Players
		teamSlug = c.PathParam("teamSlug")
	)
	log.Println(teamSlug)

	if team, err = t.TeamRepo.GetOneBySlug(teamSlug); err != nil {
		return err
	}
	if players, err = t.PlayerRepo.GetAllByTeamSlug(teamSlug); err != nil {
		return fmt.Errorf("error fetching players for %s: %s", teamSlug, err)
	}
	log.Println(team, players)
	return view.TeamRoster(c, team, players).Render(c.Request().Context(), c.Response().Writer)
}
