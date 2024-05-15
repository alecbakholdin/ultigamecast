package service

import (
	"fmt"
	"sync"
	"ultigamecast/pbmodels"
	"ultigamecast/repository"

	"github.com/google/uuid"
	"github.com/pocketbase/pocketbase/daos"
)

type LiveGames struct {
	gameRepo        *repository.Game
	eventRepo       *repository.Events
	transactionRepo *repository.Transaction

	gamesMut   sync.RWMutex
	games      []*liveGameContainer
	addGame    chan *liveGameContainer
	deleteGame chan string
}

func NewLiveGames(g *repository.Game, e *repository.Events, t *repository.Transaction) *LiveGames {
	return &LiveGames{
		gameRepo:        g,
		eventRepo:       e,
		transactionRepo: t,
		games:           make([]*liveGameContainer, 0),
		addGame:         make(chan *liveGameContainer, 5),
		deleteGame:      make(chan string, 5),
	}
}

func (l *LiveGames) StartLiveGame(gameId string) (*liveGameContainer, error) {
	l.gamesMut.Lock()
	defer l.gamesMut.Unlock()

	if lg := l.findLiveGame(gameId); lg != nil {
		return lg, nil
	}

	game, err := l.gameRepo.GetOneById(gameId)
	if err != nil {
		return nil, fmt.Errorf("error fetching game %s: %s", gameId, err)
	}
	events, err := l.eventRepo.GetAllByGame(gameId)
	if err != nil {
		return nil, fmt.Errorf("error fetching game events for game %s: %s", gameId, err)
	}

	lg := &liveGameContainer{
		mut:           &sync.RWMutex{},
		gameUpdates:   make(chan *LiveGameUpdate, 100),
		eventUpdates:  make(chan *LiveEventUpdate, 100),
		Game:          game,
		Subscriptions: make([]*LiveGameSubscription, 0),
		Events:        events,
		pointType:     pbmodels.EventsPointTypeO,
	}
	l.games = append(l.games, lg)
	go l.notifySubscribers(lg)
	return lg, nil
}

func (l *LiveGames) RemoveLiveGame(gameId string) {
	l.gamesMut.Lock()
	defer l.gamesMut.Unlock()

	for i, lg := range l.games {
		if lg.Game.Id != gameId {
			continue
		}

		lg.mut.Lock()
		for _, s := range lg.Subscriptions {
			s.game = nil
			l.Unsubscribe(s)
		}
		l.games = append(l.games[:i], l.games[i+1:]...)
		lg.mut.Unlock()
		break
	}
}

func (l *LiveGames) Subscribe(gameId string) (*LiveGameSubscription, error) {
	liveGame, err := l.StartLiveGame(gameId)
	if err != nil {
		return nil, fmt.Errorf("error starting live game: %s", err)
	}
	liveGame.mut.Lock()
	defer liveGame.mut.Unlock()

	sub := &LiveGameSubscription{
		GameId:         gameId,
		SubscriptionId: uuid.NewString(),
		GameUpdates:    make(chan *LiveGameUpdate, 10),
		EventUpdates:   make(chan *LiveEventUpdate, 10),
		Close:          make(chan int),
		game:           liveGame,
	}
	liveGame.Subscriptions = append(liveGame.Subscriptions, sub)
	return sub, nil
}

func (l *LiveGames) Unsubscribe(sub *LiveGameSubscription) {
	sub.Close <- 1
	if sub.game == nil {
		return
	}

	sub.game.mut.Lock()
	defer sub.game.mut.Unlock()

	for i, s := range sub.game.Subscriptions {
		if s.SubscriptionId == sub.SubscriptionId {
			sub.game.Subscriptions = append(sub.game.Subscriptions[:i], sub.game.Subscriptions[i+1:]...)
			break
		}
	}
}

func (l *LiveGames) StartingLine(gameId string, playerIds []string) error {
	if len(playerIds) != 7 {
		return fmt.Errorf("expected 7 players but found %d", len(playerIds))
	}
	err := l.genericUpdate(gameId, pbmodels.GamesLiveStatusBeforePoint, func(g *liveGameContainer) []*pbmodels.Events {
		events := make([]*pbmodels.Events, 7)
		for i, p := range playerIds {
			events[i] = g.newPlayerEvent(pbmodels.EventsTypeStartingLine, p, "")
		}
		return events
	})
	if err != nil {
		return fmt.Errorf("error starting point for %s: %s", gameId, err)
	}
	return nil
}

func (l *LiveGames) TeamPlayerEvent(gameId string, eventType pbmodels.EventsType, playerId, message string) error {
	err := l.genericUpdate(gameId, pbmodels.GamesLiveStatusTeamPossession, func(g *liveGameContainer) []*pbmodels.Events {
		return []*pbmodels.Events{g.newPlayerEvent(eventType, playerId, message)}
	})
	if err != nil {
		return fmt.Errorf("error creating event %s for player %s in game %s: %s", eventType, playerId, gameId, err)
	}
	return nil
}

func (l *LiveGames) OpponentPlayerEvent(gameId string, eventType pbmodels.EventsType, message string) error {
	err := l.genericUpdate(gameId, pbmodels.GamesLiveStatusTeamPossession, func(g *liveGameContainer) []*pbmodels.Events {
		return []*pbmodels.Events{g.newOpponentEvent(eventType, message)}
	})
	if err != nil {
		return fmt.Errorf("error creating event %s for opponent in game %s: %s", eventType, gameId, err)
	}
	return nil
}

func (l *LiveGames) TeamScored(gameId string, assistPlayerId, goalPlayerId string, message string) error {
	err := l.genericUpdate(gameId, pbmodels.GamesLiveStatusTeamPossession, func(g *liveGameContainer) []*pbmodels.Events {
		events := make([]*pbmodels.Events, 0)
		if assistPlayerId != "" {
			events = append(events, g.newPlayerEvent(pbmodels.EventsTypeAssist, assistPlayerId, message))
		}
		return append(events, g.newPlayerEvent(pbmodels.EventsTypeGoal, goalPlayerId, message))
	})
	if err != nil {
		return fmt.Errorf("error creating goal for [%s] [%s] in game %s with message %s: %s", assistPlayerId, goalPlayerId, gameId, message, err)
	}
	return nil
}

func (l *LiveGames) OpponentScored(gameId, message string) error {
	err := l.genericUpdate(gameId, pbmodels.GamesLiveStatusOpponentPossession, func(g *liveGameContainer) []*pbmodels.Events {
		return []*pbmodels.Events{g.newOpponentEvent(pbmodels.EventsTypeGoal, message)}
	})
	if err != nil {
		return fmt.Errorf("error creating opponent goal for in game %s with message %s: %s", gameId, message, err)
	}
	return nil
}

func (l *LiveGames) genericUpdate(gameId string, expectedStatus pbmodels.GamesLiveStatus, getEvents func(g *liveGameContainer) []*pbmodels.Events) error {
	game := l.lockedFindLiveGame(gameId)
	if game == nil {
		fmt.Println("No game found")
		return nil
	}

	game.mut.Lock()
	defer game.mut.Unlock()

	events := getEvents(game)
	fmt.Println(events)
	err := l.transactionRepo.Run(func(txDao *daos.Dao) error {
		g := game.Game.Copy()
		g.Id = game.Game.Id
		updatedFields := []string{}
		for _, e := range events {
			updatedFields = append(updatedFields, l.applyEvent(g, e, false)...)
			if err := l.eventRepo.CreateDao(txDao, e); err != nil {
				return fmt.Errorf("error creating event %v: %s", e, err)
			}
		}
		fmt.Println(g, g.OpponentScore)
		if len(updatedFields) == 0 {
			return nil
		} else if err := l.gameRepo.UpdateDao(txDao, g, updatedFields...); err != nil {
			return fmt.Errorf("error updating game: %s", err)
		}
		return nil
	})
	fmt.Println("err", err)
	if err != nil {
		return err
	}

	updatedFields := []string{}
	for _, e := range events {
		updatedFields = append(updatedFields, l.applyEvent(game.Game, e, false)...)
		game.eventUpdates <- &LiveEventUpdate{
			IsNew: true,
			Event: e,
		}
	}
	fmt.Println(events[0], updatedFields)
	if len(updatedFields) > 0 {
		game.gameUpdates <- &LiveGameUpdate{
			Fields: updatedFields,
			Game:   game.Game.Copy(),
		}
	}

	return nil
}

// updates the Games model based on the provided events, reverses if undoing.
// Returns the updated fields
func (l *LiveGames) applyEvent(g *pbmodels.Games, event *pbmodels.Events, undo bool) []string {
	updatedFields := []string{}
	switch event.Type {
	case pbmodels.EventsTypeGoal:
		if event.IsOpponent {
			updatedFields = append(updatedFields, "LiveStatus", "OpponentScore")
			if undo {
				g.OpponentScore -= 1
				g.LiveStatus = pbmodels.GamesLiveStatusOpponentPossession
			} else {
				g.OpponentScore += 1
				g.LiveStatus = pbmodels.GamesLiveStatusBeforePoint
			}
		} else {
			updatedFields = append(updatedFields, "LiveStatus", "TeamScore")
			if undo {
				g.TeamScore -= 1
				g.LiveStatus = pbmodels.GamesLiveStatusTeamPossession
			} else {
				g.TeamScore += 1
				g.LiveStatus = pbmodels.GamesLiveStatusBeforePoint
			}
		}
	case pbmodels.EventsTypeBlock:
	case pbmodels.EventsTypeDrop:
	case pbmodels.EventsTypeTurn:
		updatedFields = append(updatedFields, "LiveStatus")
		g.LiveStatus = switchPossession(g.LiveStatus)
	case pbmodels.EventsTypeHalftime:
		updatedFields = append(updatedFields, "LiveStatus")
		if undo {
			g.LiveStatus = pbmodels.GamesLiveStatusBeforePoint
		} else {
			g.LiveStatus = pbmodels.GamesLiveStatusHalftime
		}
	case pbmodels.EventsTypeGameEnd:
		updatedFields = append(updatedFields, "LiveStatus")
		if undo {
			g.LiveStatus = pbmodels.GamesLiveStatusBeforePoint
		} else {
			g.LiveStatus = pbmodels.GamesLiveStatusFinal
		}
	}
	return updatedFields
}

func switchPossession(liveStatus pbmodels.GamesLiveStatus) pbmodels.GamesLiveStatus {
	if liveStatus == pbmodels.GamesLiveStatusOpponentPossession {
		return pbmodels.GamesLiveStatusTeamPossession
	} else if liveStatus == pbmodels.GamesLiveStatusTeamPossession {
		return pbmodels.GamesLiveStatusOpponentPossession
	}
	return liveStatus
}

func (l *LiveGames) notifySubscribers(game *liveGameContainer) {
	for {
		select {
		case gu := <-game.gameUpdates:
			for _, s := range game.Subscriptions {
				s.GameUpdates <- gu
			}
		case eu := <-game.eventUpdates:
			for _, s := range game.Subscriptions {
				s.EventUpdates <- eu
			}
		}
	}
}

func (l *LiveGames) lockedFindLiveGame(gameId string) *liveGameContainer {
	l.gamesMut.Lock()
	defer l.gamesMut.Unlock()

	return l.findLiveGame(gameId)
}

// not thread-safe. Lock before this
func (l *LiveGames) findLiveGame(gameId string) *liveGameContainer {
	for _, game := range l.games {
		if game.Game.Id == gameId {
			return game
		}
	}
	return nil
}
