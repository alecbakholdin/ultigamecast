package handlers

import (
	"net/http"
	"ultigamecast/models"
	view_team "ultigamecast/view/team"
)

type Team struct {
	t TeamService
}

type TeamService interface {

}

func NewTeam(t TeamService) *Team {
	return &Team{
		t: t,
	}
}

func (t *Team) GetTeams(w http.ResponseWriter, r *http.Request) {
	view_team.TeamsPage([]*models.Team{}).Render(r.Context(), w)
}

func (t *Team) PostTeams(w http.ResponseWriter, r *http.Request) {

}
