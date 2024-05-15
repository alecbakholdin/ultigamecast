package service

import (
	"fmt"
	"ultigamecast/pbmodels"
	"ultigamecast/repository"
)

type Events struct {
	eventsRepo *repository.Events
}

func NewEvents(e *repository.Events) *Events {
	return &Events{
		eventsRepo: e,
	}
}

func (e *Events) GetAllByGame(gameId string) ([]*pbmodels.Events, error) {
	events, err := e.eventsRepo.GetAllByGame(gameId)
	if err != nil {
		return nil, fmt.Errorf("error fetching events for %s: %s", events, err)
	}
	return events, nil
}
