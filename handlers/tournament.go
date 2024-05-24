package handlers

import "net/http"

type Tournament struct {
	t TournamentService
}

type TournamentService interface {
	
}

func NewTournament(t TournamentService) *Tournament {
	return &Tournament{
		t: t,
	}
}

func (t *Tournament) GetTournaments(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

func (t *Tournament) GetTournament(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

func (t *Tournament) PostTournaments(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

func (t *Tournament) PutTournament(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented", http.StatusNotImplemented)
}