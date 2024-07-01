package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"
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
		return
	}
	defer c.CloseNow()

	// set up context
	if !ctxvar.IsAdmin(ctx) {
		ctx = c.CloseRead(ctx)
	}
	ctx, cancel = context.WithCancel(ctx)
	*r = *r.WithContext(ctx)

	// subscribe to game events
	sub, err := g.event.Subscribe(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "unexpected error subscribing to game", "err", err)
		http.Error(w, "unexpected error", http.StatusInternalServerError)
		return
	}
	defer g.event.Unsubscribe(ctx, sub.Id)

	//admin := ctxvar.IsAdmin(r.Context())
	go func() {
		for {
			msgType, bytes, err := c.Read(ctx)
			if websocket.CloseStatus(err) != -1 {
				slog.InfoContext(ctx, "closing websocket")
				fmt.Println(c.Close(websocket.StatusNormalClosure, ""))
				return
			} else if err != nil {
				slog.ErrorContext(ctx, "error reading from websocket", "err", err)
			} else if msgType != websocket.MessageText {
				slog.ErrorContext(ctx, "unexpected message type", "message type", msgType)
			} else if hMsg, err := parseMessage(bytes); err != nil {
				slog.ErrorContext(ctx, "error converting to htmxWsMessage", "err", err)
			} else {
				wsMessages <- &hMsg
			}
		}
	}()
	isAdmin := ctxvar.IsAdmin(ctx)

	ticker := time.NewTicker(time.Second * 15)
	for {
		select {
		case wsMessage := <-wsMessages:
			if isAdmin {
				fmt.Println(*wsMessage)
			}
		case u := <-sub.EventChan:
			if isAdmin {
				g.adminEventUpdate(ctx, u)
			}
			g.eventUpdate(ctx, u)
		case <-ticker.C:
			slog.InfoContext(ctx, "pinging")
			if err := c.Ping(ctx); err != nil {
				slog.ErrorContext(ctx, "client unresponsive", "err", err)
				cancel()
			}
		case <-ctx.Done():
			slog.InfoContext(ctx, "returning")
			return
		}
	}
}

func parseMessage(data []byte) (HtmxWsMessage, error) {
	m := make(map[string]any)
	msg := HtmxWsMessage{}
	if err := json.Unmarshal(data, &m); err != nil {
		return msg, err
	}
	headers, ok := m["HEADERS"].(map[string]any)
	if !ok {
		return msg, fmt.Errorf("missing headers in htmx message: %s", string(data))
	}
	msg.Trigger = headers["HX-Trigger"].(string)
	msg.Target = headers["HX-Target"].(string)
	msg.CurrentURL = headers["HX-Current-URL"].(string)
	delete(m, "HEADERS")
	msg.Data = m
	return msg, nil
}

func (g *Game) readAdminUpdates(ctx context.Context, c *websocket.Conn) chan *HtmxWsMessage {
	wsMessages := make(chan *HtmxWsMessage, 5)
	go func() {
		for {
			msgType, bytes, err := c.Read(ctx)
			if websocket.CloseStatus(err) != -1 {
				slog.InfoContext(ctx, "closing websocket")
				fmt.Println(c.Close(websocket.StatusNormalClosure, ""))
			} else if err != nil {
				slog.ErrorContext(ctx, "error reading from websocket", "err", err)
			} else if msgType != websocket.MessageText {
				slog.ErrorContext(ctx, "unexpected message type", "message type", msgType)
			} else if hMsg, err := parseMessage(bytes); err != nil {
				slog.ErrorContext(ctx, "error converting to htmxWsMessage", "err", err)
			} else {
				wsMessages <- &hMsg
			}
		}
	}()
	return wsMessages
}

func (g *Game) adminEventUpdate(ctx context.Context, e *service.EventUpdate) {
}

func (g *Game) eventUpdate(ctx context.Context, e *service.EventUpdate) {
}
