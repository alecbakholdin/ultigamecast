package handlers

import (
	"context"
	"net/http"
	"ultigamecast/internal/app/handlers/htmx"
	"ultigamecast/internal/models"
	view_team_schedule "ultigamecast/web/view/teams/schedule"
)

type TeamScheduleService interface {
	GetSchedule(ctx context.Context) ([]models.TournamentSummary, error)
	CreateTournament(ctx context.Context, name, dates string) (*models.TournamentSummary, error)
}

type TeamSchedule struct {
	schedule TeamScheduleService
}

func NewTeamSchedule(schedule TeamScheduleService) *TeamSchedule {
	return &TeamSchedule{
		schedule: schedule,
	}
}

func (t *TeamSchedule) Get(w http.ResponseWriter, r *http.Request) {
	if summaries, err := t.schedule.GetSchedule(r.Context()); err != nil {
		http.Error(w, "unexpected error", http.StatusInternalServerError)
	} else {
		view_team_schedule.Schedule(summaries).Render(r.Context(), w)
	}
}

func (t *TeamSchedule) Post(w http.ResponseWriter, r *http.Request) {
	dto := &view_team_schedule.CreateTournamentDTO{
		Name:  r.FormValue("name"),
		Dates: r.FormValue("dates"),
	}
	if !dto.Validate(dto) {
		htmx.HxCreateDatepicker(w, view_team_schedule.CreateTournamentDatesId)
		view_team_schedule.CreateTournamentForm(dto).Render(r.Context(), w)
		return
	}
	if tournament, err := t.schedule.CreateTournament(r.Context(), dto.Name, dto.Dates); err != nil {
		dto.AddFormError("unexpected error")
		htmx.HxCreateDatepicker(w, view_team_schedule.CreateTournamentDatesId)
		view_team_schedule.CreateTournamentForm(dto).Render(r.Context(), w)
		return
	} else {
		htmx.HxCloseModal(w)
		htmx.HxRetargetSwap(w, "#schedule-list", "beforeend")
		view_team_schedule.ScheduleRow(tournament).Render(r.Context(), w)
	}
}

func (t *TeamSchedule) GetModal(w http.ResponseWriter, r *http.Request) {
	htmx.HxOpenModal(w)
	htmx.HxCreateDatepicker(w, view_team_schedule.CreateTournamentDatesId)
	view_team_schedule.TournamentModal().Render(r.Context(), w)
}
