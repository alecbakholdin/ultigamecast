package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"sync"
	"ultigamecast/internal/assert"
	"ultigamecast/internal/ctxvar"
	"ultigamecast/internal/models"

	"github.com/google/uuid"
)

type Event struct {
	q  *models.Queries
	db *sql.DB

	mut  sync.RWMutex
	subs []*EventSubscription
}

type EventSubscription struct {
	Id        string
	GameId    int64
	EventChan chan *EventUpdate
}

type EventUpdate struct {
	Id string
	Event []models.Event
	Game  *models.Game
}

func NewEvent(db *sql.DB) *Event {
	return &Event{
		q:  models.New(db),
		db: db,
	}
}

func (e *Event) GameEvents(ctx context.Context) ([]models.Event, error) {
	game := ctxvar.GetGame(ctx)
	assert.That(game != nil, "game was nil")
	events, err := e.q.ListGameEvents(ctx, game.ID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, convertAndLogSqlError(ctx, "error fetching game events", err)
	}
	return events, nil
}

func (e *Event) Subscribe(ctx context.Context) (*EventSubscription, error) {
	game := ctxvar.GetGame(ctx)
	assert.That(game != nil, "game should not be nil when subscribing")
	e.mut.Lock()
	defer e.mut.Unlock()

	if game.ScheduleStatus != models.GameScheduleStatusLive {
		return nil, ErrGameNotLive
	}

	sub := EventSubscription{
		Id:        uuid.NewString(),
		GameId:    game.ID,
		EventChan: make(chan *EventUpdate, 10),
	}
	e.subs = append(e.subs, &sub)
	return &sub, nil
}

func (e *Event) Unsubscribe(ctx context.Context, id string) {
	e.mut.Lock()
	defer e.mut.Unlock()

	for i, s := range e.subs {
		if s.Id == id {
			e.subs = append(e.subs[0:i], e.subs[i+1:]...)
		}
	}
}

func (e *Event) ClearGameSubscriptions(ctx context.Context) {
	game := ctxvar.GetGame(ctx)
	assert.That(game != nil, "game cannot be nil when clearing game subscriptions")

	e.mut.Lock()
	defer e.mut.Unlock()

	for i, s := range e.subs {
		if s.GameId == game.ID {
			e.subs = append(e.subs[0:i], e.subs[i+1:]...)
		}
	}
}

func (e *Event) StartingLine(ctx context.Context, playerSlugs ...string) error {
	team := ctxvar.GetTeam(ctx)
	assert.That(team != nil, "team is nil")
	players, err := e.q.ListTeamPlayersBySlug(ctx, models.ListTeamPlayersBySlugParams{TeamId: team.ID, PlayerSlugs: playerSlugs})
	if err != nil {
		return convertAndLogSqlError(ctx, "error fetching players", err)
	} else if len(players) != len(playerSlugs) {
		slog.ErrorContext(ctx, fmt.Sprintf("expected %d players but got %d: %v", len(playerSlugs), len(players), players))
		return ErrNotFound
	}
	activePlayers := ""
	events := make([]models.CreateEventParams, len(playerSlugs))
	for i, p := range players {
		if i != 0 {
			activePlayers += ","
		}
		activePlayers += strconv.FormatInt(p.ID, 10)
		events[i] = models.CreateEventParams{Type: models.EventTypeStartingLine, Player: sql.NullInt64{Int64: p.ID, Valid: true}}
	}
	return e.applyChanges(ctx,
		models.GameLiveStatusPrePoint,
		func(ulgp *models.UpdateLiveGameParams) {
			ulgp.ActivePlayers = sql.NullString{String: activePlayers, Valid: true}
		},
		events...,
	)
}

func (e *Event) StartPointOffense(ctx context.Context) error {
	game := ctxvar.GetGame(ctx)
	assert.That(game != nil, "game cannot be nil")
	if g, err := e.q.GetGameById(ctx, game.ID); err != nil {
		return convertAndLogSqlError(ctx, "error fetching game", err)
	} else if strings.Count(g.ActivePlayers.String, ",") != 6 {
		return ErrLineNotReady
	}
	assert.That(game != nil, "game cannot be nil")
	return e.advanceState(ctx, models.EventTypePointStart, models.GameLiveStatusPrePoint, models.GameLiveStatusTeamPossession)
}

func (e *Event) StartPointDefense(ctx context.Context) error {
	game := ctxvar.GetGame(ctx)
	assert.That(game != nil, "game cannot be nil")
	if g, err := e.q.GetGameById(ctx, game.ID); err != nil {
		return convertAndLogSqlError(ctx, "error fetching game", err)
	} else if strings.Count(g.ActivePlayers.String, ",") != 6 {
		return ErrLineNotReady
	}
	return e.advanceState(ctx, models.EventTypePointStart, models.GameLiveStatusPrePoint, models.GameLiveStatusOpponentPossession)
}

func (e *Event) TimeoutPoint(ctx context.Context) error {
	return e.advanceState(ctx, models.EventTypeTimeout, models.GameLiveStatusTeamPossession, models.GameLiveStatusTeamTimeout)
}

func (e *Event) TimeoutEndPoint(ctx context.Context) error {
	return e.advanceState(ctx, models.EventTypeTimeoutEnd, models.GameLiveStatusTeamTimeout, models.GameLiveStatusTeamPossession)
}

func (e *Event) OpponentTimeoutPoint(ctx context.Context) error {
	return e.advanceState(ctx, models.EventTypeOpponentTimeout, models.GameLiveStatusOpponentPossession, models.GameLiveStatusOpponentTimeout)
}

func (e *Event) OpponentTimeoutEndPoint(ctx context.Context) error {
	return e.advanceState(ctx, models.EventTypeOpponentTimeoutEnd, models.GameLiveStatusOpponentTimeout, models.GameLiveStatusOpponentPossession)
}

func (e *Event) Timeout(ctx context.Context) error {
	return e.advanceState(ctx, models.EventTypeTimeout, models.GameLiveStatusPrePoint, models.GameLiveStatusPrePoint)
}

func (e *Event) TimeoutEnd(ctx context.Context) error {
	return e.advanceState(ctx, models.EventTypeTimeoutEnd, models.GameLiveStatusPrePoint, models.GameLiveStatusPrePoint)
}

func (e *Event) OpponentTimeout(ctx context.Context) error {
	return e.advanceState(ctx, models.EventTypeOpponentTimeout, models.GameLiveStatusPrePoint, models.GameLiveStatusOpponentTimeout)
}

func (e *Event) OpponentTimeoutEnd(ctx context.Context) error {
	return e.advanceState(ctx, models.EventTypeOpponentTimeoutEnd, models.GameLiveStatusOpponentTimeout, models.GameLiveStatusPrePoint)
}

func (e *Event) Touch(ctx context.Context, playerSlug string) error {
	return e.playerStat(ctx, models.EventTypeTouch, models.GameLiveStatusTeamPossession, models.GameLiveStatusTeamPossession, playerSlug)
}

func (e *Event) Block(ctx context.Context, playerSlug string) error {
	return e.playerStat(ctx, models.EventTypeBlock, models.GameLiveStatusOpponentPossession, models.GameLiveStatusTeamPossession, playerSlug)
}

func (e *Event) Turn(ctx context.Context, playerSlug string) error {
	return e.playerStat(ctx, models.EventTypeTurn, models.GameLiveStatusTeamPossession, models.GameLiveStatusOpponentPossession, playerSlug)
}

func (e *Event) Drop(ctx context.Context, playerSlug string) error {
	return e.playerStat(ctx, models.EventTypeDrop, models.GameLiveStatusTeamPossession, models.GameLiveStatusOpponentPossession, playerSlug)
}

func (e *Event) OpponentBlock(ctx context.Context) error {
	return e.advanceState(ctx, models.EventTypeOpponentBlock, models.GameLiveStatusTeamPossession, models.GameLiveStatusOpponentPossession)
}

func (e *Event) OpponentTurn(ctx context.Context) error {
	return e.advanceState(ctx, models.EventTypeOpponentTurn, models.GameLiveStatusOpponentPossession, models.GameLiveStatusTeamPossession)
}

func (e *Event) OpponentDrop(ctx context.Context) error {
	return e.advanceState(ctx, models.EventTypeOpponentDrop, models.GameLiveStatusOpponentPossession, models.GameLiveStatusTeamPossession)
}

func (e *Event) Goal(ctx context.Context, goalPlayerSlug, assistPlayerSlug string) error {
	team := ctxvar.GetTeam(ctx)
	assert.That(team != nil, "team cannot be null when calling Goal")
	assert.That(goalPlayerSlug != "", "goal player cannot be null when calling Goal")
	assert.That(assistPlayerSlug != "", "assist player cannot be null when calling Assist")
	goalPlayer, err := e.q.GetPlayer(ctx, models.GetPlayerParams{TeamId: team.ID, Slug: goalPlayerSlug})
	if err != nil {
		return convertAndLogSqlError(ctx, "error getting goalPlayer", err)
	}
	assistPlayer, err := e.q.GetPlayer(ctx, models.GetPlayerParams{TeamId: team.ID, Slug: assistPlayerSlug})
	if err != nil {
		return convertAndLogSqlError(ctx, "error getting assistPlayer", err)
	}
	return e.applyChanges(ctx,
		models.GameLiveStatusTeamPossession,
		func(ulgp *models.UpdateLiveGameParams) {
			ulgp.TeamScore++
			ulgp.LiveStatus = models.GameLiveStatusPrePoint
			ulgp.ActivePlayers = sql.NullString{Valid: false}
		},
		models.CreateEventParams{Type: models.EventTypeGoal, Player: sql.NullInt64{Int64: goalPlayer.ID, Valid: true}},
		models.CreateEventParams{Type: models.EventTypeAssist, Player: sql.NullInt64{Int64: assistPlayer.ID, Valid: true}},
	)
}

func (e *Event) OpponentGoal(ctx context.Context) error {
	return e.applyChanges(ctx,
		models.GameLiveStatusOpponentPossession,
		func(ulgp *models.UpdateLiveGameParams) {
			ulgp.OpponentScore++
			ulgp.LiveStatus = models.GameLiveStatusPrePoint
			ulgp.ActivePlayers = sql.NullString{Valid: false}
		},
		models.CreateEventParams{Type: models.EventTypeOpponentGoal},
	)
}

func (e *Event) playerStat(ctx context.Context, eventType models.EventType, current, next models.GameLiveStatus, playerSlug string) error {
	team := ctxvar.GetTeam(ctx)
	assert.That(team != nil, "team cannot be nil for player event "+string(eventType))
	player, err := e.q.GetPlayer(ctx, models.GetPlayerParams{TeamId: team.ID, Slug: playerSlug})
	if err != nil {
		return convertAndLogSqlError(ctx, "error finding team player %s %s", err)
	}
	return e.applyChanges(ctx,
		current,
		func(ulgp *models.UpdateLiveGameParams) {
			ulgp.LiveStatus = next
		},
		models.CreateEventParams{
			Type:   eventType,
			Player: sql.NullInt64{Int64: player.ID, Valid: true},
		},
	)
}

func (e *Event) advanceState(ctx context.Context, eventType models.EventType, current, next models.GameLiveStatus) error {
	return e.applyChanges(ctx,
		current,
		func(ulgp *models.UpdateLiveGameParams) {
			ulgp.LiveStatus = next
		},
		models.CreateEventParams{Type: eventType},
	)
}

func (e *Event) applyChanges(ctx context.Context, expectedState models.GameLiveStatus, modifierFn func(*models.UpdateLiveGameParams), events ...models.CreateEventParams) error {
	game := ctxvar.GetGame(ctx)
	assert.That(game != nil, "event changes cannot be applied to nil game")
	tx, err := e.db.Begin()
	if err != nil {
		return convertAndLogSqlError(ctx, "error opening transaction for applying event", err)
	}
	defer func() { convertAndLogSqlError(ctx, "error rolling back transaion", tx.Rollback()) }()

	q := e.q.WithTx(tx)
	liveGameData, err := q.GetGameById(ctx, game.ID)
	if err != nil {
		return convertAndLogSqlError(ctx, "error fetching game by id", err)
	}
	if liveGameData.LiveStatus != expectedState {
		return fmt.Errorf("state %s does not match expected state %s", liveGameData.LiveStatus, expectedState)
	}

	createdEvents, err := e.createEvents(ctx, q, &liveGameData, events)
	if err != nil {
		return err
	}

	gameUpdate := &models.UpdateLiveGameParams{
		ID:            liveGameData.ID,
		TeamScore:     liveGameData.TeamScore,
		OpponentScore: liveGameData.OpponentScore,
		LiveStatus:    liveGameData.LiveStatus,
		ActivePlayers: liveGameData.ActivePlayers,
		LastEvent:     sql.NullString{String: createdEvents[len(createdEvents)-1].ID, Valid: true},
	}
	if len(createdEvents) > 1 && createdEvents[0].Batch.Valid {
		gameUpdate.LastEvent = createdEvents[0].Batch
	}
	modifierFn(gameUpdate)
	updatedGame, err := q.UpdateLiveGame(ctx, *gameUpdate)
	if err != nil {
		return convertAndLogSqlError(ctx, "error updating game", err)
	}
	if err := tx.Commit(); err != nil {
		return convertAndLogSqlError(ctx, "error committing transaction", err)
	}

	go e.notifySubscribers(&updatedGame, createdEvents)
	return nil
}

func (e *Event) createEvents(ctx context.Context, qWithTx *models.Queries, game *models.Game, events []models.CreateEventParams) ([]models.Event, error) {
	var err error
	var batchId string
	if len(events) > 1 {
		batchId = uuid.NewString()
	}
	createdEvents := make([]models.Event, len(events))
	for i, e := range events {
		e.Game = game.ID
		if len(events) > 1 {
			e.Batch = sql.NullString{String: batchId, Valid: true}
		}
		e.ID = uuid.NewString()
		e.TeamScore = game.TeamScore
		e.OpponentScore = game.OpponentScore
		e.PreviousGameState = game.LiveStatus
		e.PreviousEvent = game.LastEvent
		if createdEvents[i], err = qWithTx.CreateEvent(ctx, e); err != nil {
			return nil, convertAndLogSqlError(ctx, fmt.Sprintf("error creating event %d %s", i, e.Type), err)
		}
	}
	return createdEvents, nil
}

func (e *Event) notifySubscribers(game *models.Game, events []models.Event) {
	e.mut.Lock()
	defer e.mut.Unlock()
	assert.That(len(events) > 0, "events cannot be empty")

	update := &EventUpdate{
		Event: events,
		Game: game,
	}
	if len(events) > 1 {
		update.Id = events[0].Batch.String
	} else {
		update.Id = events[0].ID
	}

	for _, s := range e.subs {
		if game.ID == s.GameId {
			s.EventChan <- update
		}
	}
}
