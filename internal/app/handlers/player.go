package handlers

import (
	"context"
	"net/http"
	"strconv"
	"ultigamecast/internal/app/handlers/htmx"
	"ultigamecast/internal/models"
	view_players "ultigamecast/web/view/teams/players"
)

type Player struct {
	p PlayerService
}

type PlayerService interface {
	GetTeamPlayers(ctx context.Context) ([]models.Player, error)
	CreatePlayer(ctx context.Context, name string) (*models.Player, error)
	UpdatePlayerOrder(ctx context.Context, playerIds []int64) error
}

func NewPlayer(p PlayerService) *Player {
	return &Player{
		p: p,
	}
}

func (p *Player) GetPlayers(w http.ResponseWriter, r *http.Request) {
	if teamPlayers, err := p.p.GetTeamPlayers(r.Context()); err != nil {
		http.Error(w, "unexpected error", http.StatusInternalServerError)
	} else {
		view_players.TeamPlayers(teamPlayers).Render(r.Context(), w)
	}
}

func (p *Player) PostPlayers(w http.ResponseWriter, r *http.Request) {
	dto := &view_players.CreatePlayerDTO{
		Name: r.FormValue("name"),
	}
	if !dto.Validate(dto) {
		view_players.CreatePlayerForm(dto).Render(r.Context(), w)
		return
	}
	if player, err := p.p.CreatePlayer(r.Context(), dto.Name); err != nil {
		dto.AddFormError("unexpected error")
		view_players.CreatePlayerForm(dto).Render(r.Context(), w)
		return
	} else {
		htmx.HxRetarget(w, "#team_players", "beforeend")
		htmx.HxClearForm(w)
		view_players.TeamPlayerRow(player).Render(r.Context(), w)
	}
}

func (p *Player) PutPlayer(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented yet", http.StatusNotImplemented)
}

func (p *Player) DeletePlayer(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented yet", http.StatusNotImplemented)
}

func (p *Player) PostPlayersOrder(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	players := r.Form["players"]
	playerInts := make([]int64, len(players))
	for i, p := range players {
		if pInt, err := strconv.Atoi(p); err != nil {
			http.Error(w, "error parsing strings", http.StatusBadRequest)
		} else {
			playerInts[i] = int64(pInt)
		}
	}

	if err := p.p.UpdatePlayerOrder(r.Context(), playerInts); err != nil {
		http.Error(w, "unexpected error", http.StatusInternalServerError)
	}
}
