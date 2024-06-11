package handlers

import (
	"net/http"
	"ultigamecast/internal/ctxvar"
	view_tournament "ultigamecast/web/view/teams/schedule/tournament"
)

type Tournament struct {
	tournament TournamentService
}

type TournamentService interface {
}

func NewTournament(tournament TournamentService) *Tournament {
	return &Tournament{
		tournament: tournament,
	}
}

func (t *Tournament) Get(w http.ResponseWriter, r *http.Request) {
	tournament := ctxvar.GetTournament(r.Context())
	if tournament == nil {
		http.Error(w, "Tournament not found", http.StatusNotFound)
	} else {
		view_tournament.TournamentPage(tournament).Render(r.Context(), w)
	}
}

func (t *Tournament) GetEdit(w http.ResponseWriter, r *http.Request) {

}

func (t *Tournament) PutEdit(w http.ResponseWriter, r *http.Request) {

}

func (t *Tournament) GetEditCancel(w http.ResponseWriter, r *http.Request) {

}
