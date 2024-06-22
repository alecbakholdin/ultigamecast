package handlers

import (
	"net/http"
	"ultigamecast/internal/app/handlers/htmx"
	"ultigamecast/internal/app/service"
	"ultigamecast/web/view/teams/roster"
)

type TeamRoster struct {
	player *service.Player
}

func NewTeamRoster(p *service.Player) *TeamRoster {
	return &TeamRoster{
		player: p,
	}
}

func (t *TeamRoster) Get(w http.ResponseWriter, r *http.Request) {
	players, err := t.player.GetTeamPlayers(r.Context())
	if err != nil {
		http.Error(w, "unexpected error", http.StatusInternalServerError)
		return
	}
	view_team_roster.TeamRoster(players).Render(r.Context(), w)
}

func (t *TeamRoster) Post(w http.ResponseWriter, r *http.Request) {
	dto := &view_team_roster.CreatePlayerDTO {
		Name: r.FormValue("name"),
	}
	if !dto.Validate(dto) {
		view_team_roster.CreatePlayerForm(dto).Render(r.Context(), w)
		return
	}

	player, err := t.player.CreatePlayer(r.Context(), dto.Name)
	if err != nil {
		dto.AddFormError("unexpected error")
		view_team_roster.CreatePlayerForm(dto).Render(r.Context(), w)
		return
	}

	htmx.HxRetargetSwap(w, "#team_roster_list", "beforeend")
	view_team_roster.TeamRosterRow(player).Render(r.Context(), w)
}

func (t *TeamRoster) PutPlayer(w http.ResponseWriter, r *http.Request) {

}

func (t *TeamRoster) DeletePlayer(w http.ResponseWriter, r *http.Request) {

}

func (t *TeamRoster) PutOrder(w http.ResponseWriter, r *http.Request) {

}
