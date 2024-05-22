package handlers

import (
	"context"
	"errors"
	"net/http"
	"ultigamecast/app/pathvar"
	"ultigamecast/models"
	"ultigamecast/service"
	view_team "ultigamecast/view/team"
)

type Team struct {
	t TeamService
}

type TeamService interface {
	GetTeam(ctx context.Context, slug string) (*models.Team, error)
	GetTeams(ctx context.Context) ([]models.Team, error)
	
	CreateTeam(ctx context.Context, name, organization string) (*models.Team, error)
}

func NewTeam(t TeamService) *Team {
	return &Team{
		t: t,
	}
}

func (t *Team) GetTeams(w http.ResponseWriter, r *http.Request) {
	if teams, err := t.t.GetTeams(r.Context()); err != nil {
		http.Error(w, "unexpected error", http.StatusInternalServerError)
	} else {
		view_team.TeamsPage(teams).Render(r.Context(), w)
	}
}

func (t *Team) GetTeamsCreate(w http.ResponseWriter, r *http.Request) {
	hxOpenModal(w)
	view_team.CreateTeamModal().Render(r.Context(), w)
}

func (t *Team) GetTeam(w http.ResponseWriter, r *http.Request) {
	if team, err := t.t.GetTeam(r.Context(), pathvar.TeamSlug(r)); errors.Is(service.ErrNotFound, err) {
		http.Error(w, "Team doesn't exist", http.StatusNotFound)
	} else if err != nil {
		http.Error(w, "unexpected error", http.StatusInternalServerError)
	} else {
		view_team.TeamPage(team).Render(r.Context(), w)
	}
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
		hxRetarget(w, "#owned_team_list", "afterbegin")
		hxCloseModal(w)
		view_team.TeamRow(team).Render(r.Context(), w)
	}
}
