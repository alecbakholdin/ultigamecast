package handlers

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"ultigamecast/internal/app/handlers/htmx"
	"ultigamecast/internal/app/service"
	"ultigamecast/internal/assert"
	"ultigamecast/internal/ctxvar"
	"ultigamecast/internal/models"
	view_game "ultigamecast/web/view/teams/schedule/tournament/schedule/game"

	"nhooyr.io/websocket"
)

type Game struct {
	game   *service.Game
	player *service.Player
	event  *service.Event
}

func NewGame(g *service.Game, p *service.Player, e *service.Event) *Game {
	return &Game{
		game:   g,
		event:  e,
		player: p,
	}
}

func (g *Game) Get(w http.ResponseWriter, r *http.Request) {
	game := ctxvar.GetGame(r.Context())
	assert.That(game != nil, "game is nil")
	playerMap, err := g.player.GetTeamPlayerMap(r.Context())
	if err != nil {
		http.Error(w, "unexpected error", http.StatusInternalServerError)
		return
	}
	events, err := g.event.GameEvents(r.Context())
	if err != nil {
		http.Error(w, "unexpected error", http.StatusInternalServerError)
		return
	}

	view_game.GamePage(game, playerMap, events).Render(r.Context(), w)
}

func (g *Game) Put(w http.ResponseWriter, r *http.Request) {
	dto := view_game.EditGameDTO{
		Field: r.FormValue("field"),
		Value: r.FormValue("value"),
	}
	if !dto.Validate(dto) {
		return
	}
	switch dto.Field {
	case models.GameFieldScheduleStatus:
		if _, err := g.game.UpdateScheduleStatus(r.Context(), dto.Value); err != nil {
			dto.AddFormError("unexpected error updating")
		} else {
			htmx.HxRefresh(w)
		}
	default:
		panic(fmt.Sprintf("Unsupported field %s", dto.Field))
	}
}

type HtmxWsMessage struct {
	Trigger    string
	Target     string
	CurrentURL string
	Data       map[string]any
}

func (g *Game) GetWs(w http.ResponseWriter, r *http.Request) {
	var (
		c          *websocket.Conn
		err        error
		ctx        context.Context = r.Context()
		cancel     context.CancelFunc
		wsMessages = make(chan *HtmxWsMessage, 5)
	)

	// open websocket
	if c, err = websocket.Accept(w, r, nil); err != nil {
		slog.ErrorContext(r.Context(), "unexpected error accepting websocket", "err", err)
		http.Error(w, "unexpected error", http.StatusInternalServerError)
		return
	}
	defer c.CloseNow()
	wsWriter, err := c.Writer(r.Context(), websocket.MessageText)
	if err != nil {
		slog.ErrorContext(r.Context(), "unexpected error creating ws writer", "err", err)
		http.Error(w, "unexpected error", http.StatusInternalServerError)
		return
	}

	sub, err := g.eventSubscriber.Subscribe(r.Context())
	if err != nil {
		slog.ErrorContext(ctx, "unexpected error subscribing to game", "err", err)
		http.Error(w, "unexpected error", http.StatusInternalServerError)
		return
	}
	defer g.eventSubscriber.Unsubscribe(r.Context(), sub.Id)
	admin := ctxvar.IsAdmin(r.Context())

	for {
		select {
		case e := <-sub.EventChan:
			if admin {
				g.adminUpdate(r.Context(), wsWriter, e)
			}
			g.update(r.Context(), wsWriter, e)
		case <-r.Context().Done():
			return
		}
	}
}

func (g *Game) adminUpdate(ctx context.Context, w io.Writer, update *service.EventUpdate) {
	for _, e := range update.Events {
		switch e.Type {

		}
	}
}

func (g *Game) update(ctx context.Context, w io.Writer, update *service.EventUpdate) {
	for _, e := range update.Events {
		switch e.Type {

		}
	}
}
