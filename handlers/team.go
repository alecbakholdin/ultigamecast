package handlers

import (
	"context"
	"errors"
	"net/http"
	"ultigamecast/models"
	"ultigamecast/service"
	view_team "ultigamecast/view/team"
)

type Team struct {
	t TeamService
}

type TeamService interface {
	GetTeams(ctx context.Context) ([]*models.Team, error)
	CreateTeam(ctx context.Context, name, organization string) (*models.Team, error)
}

func NewTeam(t TeamService) *Team {
	return &Team{
		t: t,
	}
}

func (t *Team) GetTeams(w http.ResponseWriter, r *http.Request) {
	view_team.TeamsPage([]*models.Team{}).Render(r.Context(), w)
}

func (t *Team) GetTeamsCreate(w http.ResponseWriter, r *http.Request) {
	hxOpenModal(w)
	view_team.CreateTeamModal().Render(r.Context(), w)
}

func (t *Team) PostTeams(w http.ResponseWriter, r *http.Request) {
	dto := &view_team.CreateTeamDTO{
		Name:         r.FormValue("name"),
		Organization: r.FormValue("organization"),
	}
	if !dto.Validate(dto) {
		view_team.CreateTeamForm(dto).Render(r.Context(), w)
		return
	}

	if team, err := t.t.CreateTeam(r.Context(), dto.Name, dto.Organization); errors.Is(service.ErrTeamExists, err) {
		dto.AddFieldError("Name", "Name is already taken")
		view_team.CreateTeamForm(dto).Render(r.Context(), w)
	} else if err != nil {
		dto.AddFormError("unexpected error")
		view_team.CreateTeamForm(dto).Render(r.Context(), w)
	} else {
		hxRetarget(w, "#team_list", "afterbegin")
		hxCloseModal(w)
		view_team.TeamRow(team).Render(r.Context(), w)
	}
}
