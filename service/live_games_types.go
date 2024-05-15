package service

import (
	"sync"
	"ultigamecast/pbmodels"
)

type liveGameContainer struct {
	mut           *sync.RWMutex
	Game          *pbmodels.Games
	Events        []*pbmodels.Events
	Subscriptions []*LiveGameSubscription

	gameUpdates  chan *LiveGameUpdate
	eventUpdates chan *LiveEventUpdate

	pointType pbmodels.EventsPointType
}

func (l *liveGameContainer) newPlayerEvent(t pbmodels.EventsType, playerId, message string) *pbmodels.Events {
	e := l.newEvent(t, message)
	e.Player = playerId
	return e
}

func (l *liveGameContainer) newOpponentEvent(t pbmodels.EventsType, message string) *pbmodels.Events {
	e := l.newEvent(t, message)
	e.IsOpponent = true
	return e
}

func (l *liveGameContainer) newEvent(t pbmodels.EventsType, message string) *pbmodels.Events {
	return &pbmodels.Events{
		IsOpponent: false,
		Game:       l.Game.Id,
		PointType:  l.pointType,
		Type:       t,
		Message:    message,
	}
}

type LiveGameUpdate struct {
	Fields []string
	Game   *pbmodels.Games
}

type LiveEventUpdate struct {
	IsNew bool
	Event *pbmodels.Events
}

type LiveGameSubscription struct {
	GameId         string
	SubscriptionId string
	GameUpdates    chan *LiveGameUpdate
	EventUpdates   chan *LiveEventUpdate
	Close          chan int
	game           *liveGameContainer
}
