package handlers

import (
	"context"
	"net/http"
	"ultigamecast/models"

	view_tournament "ultigamecast/view/teams/tournaments"
)

type Tournament struct {
	t TournamentService
}

type TournamentService interface {
	CreateTournament(ctx context.Context, name string) (*models.Tournament, error)
	GetTeamTournaments(ctx context.Context) ([]models.Tournament, error)
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
