package handlers

import (
	"net/http"
	"ultigamecast/internal/assert"
	"ultigamecast/internal/ctxvar"
	view_game "ultigamecast/web/view/teams/schedule/tournament/schedule/game"
)

type Game struct {

}

func NewGame() *Game {
	return &Game{}
}

func (g *Game) Get(w http.ResponseWriter, r *http.Request) {
	game := ctxvar.GetGame(r.Context())
	assert.That(game != nil, "game is nil")
	view_game.GamePage(game).Render(r.Context(), w)
}

func (g *Game) GetWs(w http.ResponseWriter, r *http.Request) {

}