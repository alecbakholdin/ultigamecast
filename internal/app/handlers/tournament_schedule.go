package handlers

import (
	"cmp"
	"context"
	"net/http"
	"strconv"
	"ultigamecast/internal/app/handlers/htmx"
	"ultigamecast/internal/models"
	view_tournament_schedule "ultigamecast/web/view/teams/schedule/tournament/schedule"
)

type TournamentSchedule struct {
	schedule TournamentScheduleService
}

type TournamentScheduleService interface {
	GetSchedule(ctx context.Context) ([]models.Game, error)
	CreateGame(ctx context.Context, opponent, start, startTimezone string, half, soft, hard int) (*models.Game, error)
}

func NewTournamentSchedule(s TournamentScheduleService) *TournamentSchedule {
	return &TournamentSchedule{
		schedule: s,
	}
}

func (t *TournamentSchedule) Get(w http.ResponseWriter, r *http.Request) {
	games, err := t.schedule.GetSchedule(r.Context())
	if err != nil {
		http.Error(w, "unexpected error", http.StatusInternalServerError)
	} else {
		view_tournament_schedule.Schedule(games).Render(r.Context(), w)
	}
}

func (t *TournamentSchedule) Post(w http.ResponseWriter, r *http.Request) {
	var err error
	dto := &view_tournament_schedule.CreateGameDTO{
		Opponent: r.FormValue("opponent"),
		Start: r.FormValue("start"),
		StartTimezone: r.FormValue("start_timezone"),
	}
	dto.HalfCap, err = strconv.Atoi(cmp.Or(r.FormValue("half_cap"), "0"))
	if err != nil {
		dto.AddFieldError("HalfCap", err.Error())
	}
	dto.SoftCap, err = strconv.Atoi(cmp.Or(r.FormValue("soft_cap"), "0"))
	if err != nil {
		dto.AddFieldError("SoftCap", err.Error())
	}
	dto.HardCap, err = strconv.Atoi(cmp.Or(r.FormValue("hard_cap"), "0"))
	if err != nil {
		dto.AddFieldError("HardCap", err.Error())
	}
	if !dto.Validate(dto) {
		view_tournament_schedule.CreateGameForm(dto).Render(r.Context(), w)
		return
	}
	_, err = t.schedule.CreateGame(r.Context(), dto.Opponent, dto.Start, dto.StartTimezone, dto.HalfCap, dto.SoftCap, dto.HardCap)
	if err != nil {
		dto.AddFormError("unexpected error")
		view_tournament_schedule.CreateGameForm(dto).Render(r.Context(), w)
		return
	}
	htmx.HxCloseModal(w)
	htmx.HxRetargetSwap(w, "#tab-0", "innerHTML")
	games, err := t.schedule.GetSchedule(r.Context())
	if err != nil {
		dto.AddFormError("unexpected error")
		view_tournament_schedule.CreateGameForm(dto).Render(r.Context(), w)
		return
	}
	view_tournament_schedule.Schedule(games).Render(r.Context(), w)
}

func (t *TournamentSchedule) GetModal(w http.ResponseWriter, r *http.Request) {
	htmx.HxOpenModal(w)
	view_tournament_schedule.GameModal().Render(r.Context(), w)
}
