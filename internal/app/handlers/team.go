package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"ultigamecast/internal/app/service"
	"ultigamecast/internal/ctxvar"
	"ultigamecast/internal/models"
	"ultigamecast/internal/pathvar"
	view_team "ultigamecast/web/view/teams"
)

type Team struct {
	t TeamService
}

type TeamService interface {
	GetTeam(ctx context.Context, slug string) (*models.Team, error)
	GetTeams(ctx context.Context) ([]models.Team, error)

	CreateTeam(ctx context.Context, name, organization string) (*models.Team, error)
	UpdateName(ctx context.Context, name string) (*models.Team, error)
	UpdateOrganization(ctx context.Context, organization string) (*models.Team, error)
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
	view_team.TeamModal(true, &view_team.TeamFormDTO{IsFirstTeam: r.URL.Query().Get("firstTeam") != ""}).Render(r.Context(), w)
}

func (t *Team) GetTeam(w http.ResponseWriter, r *http.Request) {
	if team, err := t.t.GetTeam(r.Context(), pathvar.TeamSlug(r)); errors.Is(err, service.ErrNotFound) {
		http.Error(w, "Team doesn't exist", http.StatusNotFound)
	} else if err != nil {
		http.Error(w, "unexpected error", http.StatusInternalServerError)
	} else {
		view_team.TeamPage(team).Render(r.Context(), w)
	}
}

func (t *Team) PostTeams(w http.ResponseWriter, r *http.Request) {
	dto := &view_team.TeamFormDTO{
		IsFirstTeam:  r.FormValue("firstTeam") != "",
		Name:         r.FormValue("name"),
		Organization: r.FormValue("organization"),
	}
	if !dto.Validate(dto) {
		view_team.TeamForm(true, dto).Render(r.Context(), w)
		return
	}

	if team, err := t.t.CreateTeam(r.Context(), dto.Name, dto.Organization); errors.Is(err, service.ErrTeamExists) {
		dto.AddFieldError("Name", "Name is already taken")
		view_team.TeamForm(true, dto).Render(r.Context(), w)
	} else if err != nil {
		dto.AddFormError("unexpected error")
		view_team.TeamForm(true, dto).Render(r.Context(), w)
	} else if dto.IsFirstTeam {
		hxRefresh(w)
	} else {
		hxCloseModal(w)
		hxRetarget(w, "#owned-teams-list", "beforeend")
		view_team.TeamRow(team).Render(r.Context(), w)
	}
}

func (t *Team) GetEdit(w http.ResponseWriter, r *http.Request) {
	field := r.URL.Query().Get("field")
	switch field {
	case "Organization":
		view_team.OrganizationForm(ctxvar.GetTeam(r.Context())).Render(r.Context(), w)
	case "Name":
		view_team.NameForm(ctxvar.GetTeam(r.Context())).Render(r.Context(), w)
	default:
		panic(fmt.Sprintf("unexpected field %s", field))
	}
}

func (t *Team) PutEdit(w http.ResponseWriter, r *http.Request) {
	field := r.URL.Query().Get("field")
	switch field {
	case "Organization":
		if team, err := t.t.UpdateOrganization(r.Context(), r.FormValue("organization")); err != nil {
			view_team.OrganizationForm(ctxvar.GetTeam(r.Context())).Render(r.Context(), w)
		} else {
			view_team.Organization(team).Render(r.Context(), w)
		}
	case "Name":
		if team, err := t.t.UpdateName(r.Context(), r.FormValue("name")); err != nil {
			view_team.NameForm(ctxvar.GetTeam(r.Context())).Render(r.Context(), w)
		} else {
			hxLocation(w, ctxvar.Url(r.Context(), team))
			view_team.Name(team).Render(r.Context(), w)
		}
	default:
		panic(fmt.Sprintf("unexpected field %s", field))
	}
}

func (t *Team) GetCancelEdit(w http.ResponseWriter, r *http.Request) {
	field := r.URL.Query().Get("field")
	switch field {
	case "Organization":
		view_team.Organization(ctxvar.GetTeam(r.Context())).Render(r.Context(), w)
	case "Name":
		view_team.Name(ctxvar.GetTeam(r.Context())).Render(r.Context(), w)
	default:
		panic(fmt.Sprintf("unexpected field %s", field))
	}
}
