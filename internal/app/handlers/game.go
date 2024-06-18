package handlers

import (
	"context"
	"log/slog"
	"net/http"
	"ultigamecast/internal/app/service"
	"ultigamecast/internal/assert"
	"ultigamecast/internal/ctxvar"
	view_game "ultigamecast/web/view/teams/schedule/tournament/schedule/game"

	"nhooyr.io/websocket"
)
type Game struct {
	eventSubscriber EventSubscriber
}

type EventSubscriber interface {
	Subscribe(ctx context.Context) (*service.EventSubscription, error)
	Unsubscribe(ctx context.Context, subId string)
}

func NewGame(eventSubscriber EventSubscriber) *Game {
	return &Game{
		eventSubscriber: eventSubscriber,
	}
}

func (g *Game) Get(w http.ResponseWriter, r *http.Request) {
	game := ctxvar.GetGame(r.Context())
	assert.That(game != nil, "game is nil")
	view_game.GamePage(game).Render(r.Context(), w)
}

func (g *Game) GetWs(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Accept(w, r, nil)
	if err != nil {
		slog.ErrorContext(r.Context(), "unexpected error accepting websocket", "err", err)
		http.Error(w, "unexpected error", http.StatusInternalServerError)
		return
	}
	defer c.CloseNow()

	sub, err := g.eventSubscriber.Subscribe(r.Context())
	if err != nil {
		slog.ErrorContext(r.Context(), "unexpected error subscribing to game", "err", err)
		http.Error(w, "unexpected error", http.StatusInternalServerError)
		return
	}
	defer g.eventSubscriber.Unsubscribe(r.Context(), sub.Id)
	//admin := ctxvar.IsAdmin(r.Context())

	for {
		select {
		//case u := <- sub.EventChan: 
			
		case <- r.Context().Done():
			return
		}
	}
}

func (g *Game) adminUpdate(e service.EventUpdate) {
}

func (g *Game) update(e service.EventUpdate) {

}