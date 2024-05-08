package service

import (
	"fmt"
	"sync"
	"ultigamecast/pbmodels"
	"ultigamecast/repository"

	"github.com/google/uuid"
)

type LiveGames struct {
	gameRepo  *repository.Game
	eventRepo *repository.Events

	gamesMut sync.RWMutex
	games    map[string]*liveGameContainer
}

func NewLiveGames(g *repository.Game, e *repository.Events) *LiveGames {
	return &LiveGames{
		gameRepo:  g,
		eventRepo: e,
	}
}

type liveGameContainer struct {
	mut           *sync.RWMutex
	Game          *pbmodels.Games
	Events        *pbmodels.Events
	Subscriptions []*LiveGameSubscription
}

type LiveGameUpdate struct {
	Field string
	Game  *pbmodels.Games
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
	game           *liveGameContainer
}

func (l *LiveGames) Subscribe(gameId string) (*LiveGameSubscription, error) {
	l.gamesMut.Lock()
	defer l.gamesMut.Unlock()

	game, ok := l.games[gameId]
	if !ok {
		return nil, fmt.Errorf("game %s is not live", gameId)
	}
	game.mut.Lock()
	defer game.mut.Unlock()

	sub := &LiveGameSubscription{
		GameId:         gameId,
		SubscriptionId: uuid.NewString(),
		GameUpdates:    make(chan *LiveGameUpdate, 10),
		EventUpdates:   make(chan *LiveEventUpdate, 10),
		game:           game,
	}
	game.Subscriptions = append(game.Subscriptions, sub)
	return sub, nil
}

func (l *LiveGames) Unsubscribe(sub *LiveGameSubscription)  {
	if sub.game == nil {
		return
	}
	sub.game.mut.Lock()
	defer sub.game.mut.Unlock()

	
	return
}
