package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"ultigamecast/internal/app/handlers/htmx"
	"ultigamecast/internal/app/service"
	"ultigamecast/internal/ctxvar"
	"ultigamecast/internal/models"
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
	htmx.HxOpenModal(w)
	view_team.TeamModal(true, &view_team.TeamFormDTO{IsFirstTeam: r.URL.Query().Get("firstTeam") != ""}).Render(r.Context(), w)
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
		htmx.HxRefresh(w)
	} else {
		htmx.HxCloseModal(w)
		htmx.HxRetarget(w, "#owned-teams-list", "beforeend")
		view_team.TeamRow(team).Render(r.Context(), w)
	}
}

func (t *Team) GetTeam(w http.ResponseWriter, r *http.Request) {
	team := ctxvar.GetTeam(r.Context())
	view_team.TeamPage(team).Render(r.Context(), w)
}

func (t *Team) GetEdit(w http.ResponseWriter, r *http.Request) {
	field := r.URL.Query().Get("field")
	team := ctxvar.GetTeam(r.Context())
	switch field {
	case "Organization":
		dto := &view_team.OrganizationDTO{Organization: team.Organization.String}
		view_team.OrganizationForm(dto).Render(r.Context(), w)
	case "Name":
		dto := &view_team.NameDTO{Name: team.Name}
		view_team.NameForm(dto).Render(r.Context(), w)
	default:
		panic(fmt.Sprintf("unexpected field %s", field))
	}
}

func (t *Team) PutEdit(w http.ResponseWriter, r *http.Request) {
	field := r.URL.Query().Get("field")
	switch field {
	case "Organization":
		dto := &view_team.OrganizationDTO{Organization: r.FormValue("organization")}
		if !dto.Validate(dto) {
			view_team.OrganizationForm(dto).Render(r.Context(), w)
			return
		}
		if team, err := t.t.UpdateOrganization(r.Context(), dto.Organization); err != nil {
			dto.AddFormError("unexpected error")
			view_team.OrganizationForm(dto).Render(r.Context(), w)
		} else {
			view_team.Organization(team).Render(r.Context(), w)
		}
	case "Name":
		dto := &view_team.NameDTO{Name: r.FormValue("name")}
		if !dto.Validate(dto) {
			view_team.NameForm(dto).Render(r.Context(), w)
			return
		}
		if team, err := t.t.UpdateName(r.Context(), dto.Name); err != nil {
			dto.AddFormError("unexpected error")
			view_team.NameForm(dto).Render(r.Context(), w)
		} else {
			htmx.HxRedirect(w, ctxvar.Url(r.Context(), team))
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