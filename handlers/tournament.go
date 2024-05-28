package handlers

import (
	"context"
	"fmt"
	"net/http"
	"ultigamecast/app/ctxvar"
	"ultigamecast/models"

	view_tournament "ultigamecast/view/teams/tournaments"
)

type Tournament struct {
	t TournamentService
}

type TournamentService interface {
	GetTeamTournaments(ctx context.Context) ([]models.Tournament, error)
	CreateTournament(ctx context.Context, name string) (*models.Tournament, error)
	UpdateTournamentDates(ctx context.Context, dates string) (*models.Tournament, error)
	DateFormat() string
}

func NewTournament(t TournamentService) *Tournament {
	return &Tournament{
		t: t,
	}
}

func (t *Tournament) GetTournaments(w http.ResponseWriter, r *http.Request) {
	tournaments, err := t.t.GetTeamTournaments(r.Context())
	if err != nil {
		http.Error(w, "unexpected error", http.StatusInternalServerError)
		return
	}
	view_tournament.TeamTournaments(tournaments).Render(r.Context(), w)
}

func (t *Tournament) GetTournament(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

func (t *Tournament) GetTournamentRow(w http.ResponseWriter, r *http.Request) {
	tournament := ctxvar.GetTournament(r.Context())
	view_tournament.TournamentRow(tournament, false).Render(r.Context(), w)
}

func (t *Tournament) GetEditDate(w http.ResponseWriter, r *http.Request) {
	tournament := ctxvar.GetTournament(r.Context())
	dto := &view_tournament.EditTournamentDatesDTO{}
	if tournament.StartDate.Valid && tournament.EndDate.Valid {
		dto.Dates = fmt.Sprintf("%s - %s", tournament.StartDate.Time.Format(t.t.DateFormat()), tournament.EndDate.Time.Format(t.t.DateFormat()))
	}
	view_tournament.EditTournamentDates(dto).Render(r.Context(), w)
}

func (t *Tournament) PutEditDate(w http.ResponseWriter, r *http.Request) {
	dto := &view_tournament.EditTournamentDatesDTO{
		Dates: r.FormValue("dates"),
	}
	if !dto.Validate(dto){ 
		view_tournament.EditTournamentDates(dto).Render(r.Context(), w)
		return
	}
	tournament, err := t.t.UpdateTournamentDates(r.Context(), dto.Dates)
	if err != nil {
		dto.AddFormError("unexpected error")
		view_tournament.EditTournamentDates(dto).Render(r.Context(), w)
		return
	}
	view_tournament.TournamentDates(tournament, false).Render(r.Context(), w)
}

func (t *Tournament) PostTournaments(w http.ResponseWriter, r *http.Request) {
	dto := &view_tournament.CreateTournamentDTO{
		Name: r.FormValue("name"),
	}
	if !dto.Validate(dto) {
		view_tournament.CreateTournamentForm(dto).Render(r.Context(), w)
		return
	}
	tournament, err := t.t.CreateTournament(r.Context(), dto.Name)
	if err != nil {
		dto.AddFormError("unexpected error")
		return
	} 
	view_tournament.CreateTournamentForm(&view_tournament.CreateTournamentDTO{}).Render(r.Context(), w)
	view_tournament.NewTournamentRow(tournament).Render(r.Context(), w)
}

func (t *Tournament) PutTournament(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented", http.StatusNotImplemented)
}
