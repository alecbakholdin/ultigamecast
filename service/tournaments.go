package service

import (
	"fmt"
	"ultigamecast/pbmodels"
	"ultigamecast/repository"
)

type Tournaments struct {
	TournamentRepo *repository.Tournament
	GameRepo       *repository.Game
	TeamRepo       *repository.Team
}

func NewTournaments(to *repository.Tournament, g *repository.Game, te *repository.Team) *Tournaments {
	return &Tournaments{
		TournamentRepo: to,
		GameRepo:       g,
		TeamRepo:       te,
	}
}

func (t *Tournaments) GetTournamentsWithGamesByTeamSlug(teamSlug string) ([]*pbmodels.TournamentsWithGames, error) {
	var (
		tournaments []*pbmodels.Tournaments
		games       []*pbmodels.Games
		tg          []*pbmodels.TournamentsWithGames
		err         error
	)

	if tournaments, err = t.TournamentRepo.GetAllByTeamSlug(teamSlug); err != nil {
		return nil, fmt.Errorf("error fetching tournaments by slug: %s", err)
	} else if games, err = t.GameRepo.GetAllForTeamBySlug(teamSlug); err != nil {
		return nil, fmt.Errorf("error fetching teams by slug: %s", err)
	}

	tg = make([]*pbmodels.TournamentsWithGames, len(tournaments))

	for i, t := range tournaments {
		tg[i] = &pbmodels.TournamentsWithGames{
			Tournament: t,
			Games:      make([]*pbmodels.Games, 0),
		}
		for _, g := range games {
			if g.Tournament == t.GetId() {
				tg[i].Games = append(tg[i].Games, g)
			}
		}
	}
	return tg, nil
}
