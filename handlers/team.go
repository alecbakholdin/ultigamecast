package handlers

import (
	"net/http"
	"ultigamecast/models"
	"ultigamecast/service"
	"ultigamecast/view/team"
)

type Team struct {
	t *service.Team
}

func NewTeam(t *service.Team) *Team {
	return &Team{
		t: t,
	}
}

func (t *Team) GetTeams(w http.ResponseWriter, r *http.Request) {
	view_team.TeamsPage([]*models.Team{}).Render(r.Context(), w)
}

func (t *Team) PostTeams(w http.ResponseWriter, r *http.Request) {

}
