package handlers

import (
	"log"
	"ultigamecast/modelspb"
	"ultigamecast/repository"
	view "ultigamecast/view/team"

	"github.com/labstack/echo/v5"
)

type Team struct {
	PlayerRepo *repository.Player
	TeamRepo   *repository.Team
}

func NewTeam(t *repository.Team, p *repository.Player) *Team {
	return &Team{
		TeamRepo:   t,
		PlayerRepo: p,
	}
}

func (t *Team) Routes(e *echo.Echo) {
	group := e.Group("/team")
	group.GET("/:teamSlug", t.getTeam)
}

func (t *Team) getTeam(c echo.Context) (err error) {
	var (
		team    *modelspb.Teams
		players []*modelspb.Players
		teamSlug = c.PathParam("teamSlug")
	)

	if team, err = t.TeamRepo.GetOneBySlug(teamSlug); err != nil {
		return err
	}
	if players, err = t.PlayerRepo.GetAllForTeamBySlug(teamSlug); err != nil {
		return err
	}
	log.Println(team, players)
	return view.Team(team, players).Render(c.Request().Context(), c.Response().Writer)
}
