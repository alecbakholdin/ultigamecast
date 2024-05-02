package repository

import (
	"log"
	"slices"
	"ultigamecast/modelspb/dto"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
)

type LiveGame struct {
	subscriptions []liveGameSubscription
}

func NewLiveGame() *LiveGame {
	return &LiveGame{
		subscriptions: make([]liveGameSubscription, 0),
	}
}

type LiveGameEvent struct {
	Event   string
	GameDto *dto.Games
}

type liveGameSubscription struct {
	subscriptionId string
	context        echo.Context
	gameId         string
	EventChan      chan LiveGameEvent
}

func (l *LiveGame) Subscribe(c echo.Context, gameId string) (subscriptionId string, eventChan chan LiveGameEvent) {
	newSub := liveGameSubscription{
		gameId:         gameId,
		context:        c,
		subscriptionId: uuid.NewString(),
		EventChan:      make(chan LiveGameEvent, 10),
	}
	log.Printf("New subscription for %s\n", gameId)
	l.subscriptions = append(l.subscriptions, newSub)
	return newSub.subscriptionId, newSub.EventChan
}

func (l *LiveGame) Unsubscribe(subscriptionId string) {
	targetIdx := slices.IndexFunc(l.subscriptions, func(s liveGameSubscription) bool {
		return s.subscriptionId == subscriptionId
	})
	if targetIdx >= 0 {
		log.Printf("Unsubscribing %s for game %s\n", subscriptionId, l.subscriptions[targetIdx].gameId)
		l.subscriptions = append(l.subscriptions[0:targetIdx], l.subscriptions[targetIdx+1:]...)
	}
}
