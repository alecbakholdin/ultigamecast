package service

import (
	"fmt"
	"strings"
	"ultigamecast/pbmodels"
	"ultigamecast/repository"
	"ultigamecast/validation"

	"github.com/labstack/echo/v5"
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

func (t *Tournaments) ValidateBasicFields(c echo.Context, to *pbmodels.Tournaments) {
	if strings.TrimSpace(to.Name) == "" {
		validation.AddFieldErrorString(c, "name", "name is required")
	}

	endDt, endErr := to.GetEndDt()
	if endErr != nil {
		validation.AddFieldErrorString(c, "end", "error parsing date")
	}

	startDt, startErr := to.GetStartDt()
	if startErr != nil {
		validation.AddFieldErrorString(c, "start", "error parsing date")
	}

	if endErr == nil && startErr == nil && !endDt.IsZero() && !startDt.IsZero() && endDt.Time().Before(startDt.Time()) {
		validation.AddFieldErrorString(c, "end", "end cannot be before start")
	}
}

func (t *Tournaments) ValidateSlugChange(c echo.Context, teamSlug, oldSlug, newSlug string) error {
	var (
		oldT *pbmodels.Tournaments
		newT *pbmodels.Tournaments
		err  error
	)
	if newSlug == "" {
		validation.AddFieldErrorString(c, "name", "name cannot be empty")
		return nil
	}
	if oldSlug != "" {
		if oldT, err = t.TournamentRepo.GetOneBySlug(teamSlug, oldSlug); err != nil {
			return fmt.Errorf("error fetching tournament with slugs [%s] [%s]: %s", teamSlug, oldSlug, err)
		}
	}
	if newT, err = t.TournamentRepo.GetOneBySlug(teamSlug, newSlug); err != nil && !repository.IsNotFound(err) {
		return fmt.Errorf("error fetching tournament with slugs [%s] [%s]: %s", teamSlug, newSlug, err)
	}

	slugAvailable := newT == nil
	isUnchanged := oldT != nil && newT != nil && newT.Id == oldT.Id
	if !(isUnchanged || slugAvailable) {
		validation.AddFieldErrorString(c, "name", "name is already taken")
	}
	return nil
}
