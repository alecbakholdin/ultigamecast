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

func (t *Tournaments) GetOneBySlug(teamSlug, tournamentSlug string) (*pbmodels.Tournaments, error) {
	if tournament, err := t.TournamentRepo.GetOneBySlug(teamSlug, tournamentSlug); err != nil {
		return nil, fmt.Errorf("error fetching tournament [%s] [%s]: %s", teamSlug, tournamentSlug, err)
	} else {
		return tournament, nil
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
	} else if games, err = t.GameRepo.GetAllByTeamSlug(teamSlug); err != nil {
		return nil, fmt.Errorf("error fetching teams by slug: %s", err)
	}

	tg = make([]*pbmodels.TournamentsWithGames, len(tournaments))

	for i, t := range tournaments {
		tg[i] = &pbmodels.TournamentsWithGames{
			Tournament: t,
			Games:      make([]*pbmodels.Games, 0),
		}
		for _, g := range games {
			fmt.Println("id", g)
			if g.Tournament == t.GetId() {
				tg[i].Games = append(tg[i].Games, g)
			}
		}
	}
	return tg, nil
}

func (t *Tournaments) Create(teamSlug string, tournament *pbmodels.Tournaments) (err error) {
	if tournament.Slug == "" {
		return fmt.Errorf("cannot have empty slug")
	}
	team, err := t.TeamRepo.FindOneBySlug(teamSlug)
	if err != nil {
		return fmt.Errorf("error creating tournament [%s] [%s]: %s", teamSlug, tournament.Slug, err)
	}
	tournament.Team = team.Id
	if err = t.TournamentRepo.Create(tournament); err != nil {
		return fmt.Errorf("error creating tournament [%s] [%s]: %s", teamSlug, tournament.Slug, err)
	}
	return nil
}

func (t *Tournaments) UpdateBySlug(teamSlug, tournamentSlug string, tournament *pbmodels.Tournaments, attrs ...string) (err error) {
	var to *pbmodels.Tournaments
	if to, err = t.TournamentRepo.GetOneBySlug(teamSlug, tournamentSlug); err != nil {
		return fmt.Errorf("error finding tournament during update of [%s] [%s]: %s", teamSlug, tournamentSlug, err)
	}
	tournament.Id = to.Id
	if err := t.TournamentRepo.Update(tournament, attrs...); err != nil {
		return fmt.Errorf("error updating tournament [%s] [%s]: %s", teamSlug, tournamentSlug, err)
	}
	return nil
}

func (t *Tournaments) ValidateBasicFields(c echo.Context, to *pbmodels.Tournaments) {
	if strings.TrimSpace(to.Name) == "" {
		validation.AddFieldErrorString(c, "name", "name is required")
	}

	var startErr, endErr error

	fmt.Println(to)
	to.End, endErr = to.GetEndDt()
	if endErr != nil {
		c.Echo().Logger.Error(fmt.Errorf("error parsing end date %s %s: %s", to.EndDatetime, to.EndTimezone, endErr))
		validation.AddFieldErrorString(c, "end", "error parsing date")
	}

	to.Start, startErr = to.GetStartDt()
	if startErr != nil {
		c.Echo().Logger.Error(fmt.Errorf("error parsing start date %s %s: %s", to.StartDatetime, to.StartTimezone, startErr))
		validation.AddFieldErrorString(c, "start", "error parsing date")
	}

	fmt.Println(to)

	if endErr == nil && startErr == nil && !to.End.IsZero() && !to.Start.IsZero() && to.End.Time().Before(to.Start.Time()) {
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

func (t *Tournaments) Delete(teamSlug, tournamentSlug string) (err error) {
	var tournament *pbmodels.Tournaments
	if tournament, err = t.TournamentRepo.GetOneBySlug(teamSlug, tournamentSlug); err != nil && !repository.IsNotFound(err) {
		return fmt.Errorf("error finding tournament while deleting [%s] [%s]: %s", teamSlug, tournamentSlug, err)
	} else if tournament == nil {
		return nil
	}

	if err = t.TournamentRepo.DeleteById(tournament.Id); err != nil {
		return fmt.Errorf("error deleting [%s] [%s]: %s", teamSlug, tournamentSlug, err)
	}

	return nil
}
