package repository

import (
	"fmt"
	"slices"
	"strings"
	"ultigamecast/modelspb"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/types"
)

type Tournament struct {
	app        core.App
	dao        *daos.Dao
	collection *models.Collection
}

func NewTournament(app core.App) *Tournament {
	return &Tournament{
		app:        app,
		dao:        app.Dao(),
		collection: mustGetCollection(app.Dao(), "tournaments"),
	}
}

func (t *Tournament) GetOneBySlug(teamSlug string, tournamentSlug string) (*modelspb.Tournaments, error) {
	record, err := t.dao.FindFirstRecordByFilter(
		t.collection.Name,
		"team.slug = {:teamSlug} && slug = {:tournamentSlug}",
		dbx.Params{
			"teamSlug":       teamSlug,
			"tournamentSlug": tournamentSlug,
		},
	)
	if err != nil {
		return nil, err
	}
	return toTournament(record), err
}

func (t *Tournament) GetOneById(id string) (*modelspb.Tournaments, error) {
	if record, err := t.dao.FindRecordById(t.collection.Id, id); err != nil {
		return nil, err
	} else {
		return toTournament(record), nil
	}
}

func (t *Tournament) ExistsBySlug(teamSlug string, tournamentSlug string) (bool, error) {
	tournament, err := t.GetOneBySlug(teamSlug, tournamentSlug)
	if IsNotFound(err) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return tournament != nil, nil
}

func (t *Tournament) GetAllByTeamSlug(slug string) ([]*modelspb.Tournaments, error) {
	records, err := t.dao.FindRecordsByFilter(
		t.collection.Name,
		"team.slug = {:teamSlug}",
		"-start",
		0,
		0,
		dbx.Params{"teamSlug": strings.ToLower(slug)},
	)
	if err != nil {
		return nil, err
	}
	return toArr(records, toTournament), nil
}

func (t *Tournament) GetAllWithGamesByTeamSlug(slug string) ([]*modelspb.TournamentWithGames, error) {
	records, err := t.dao.FindRecordsByFilter(
		t.collection.Name,
		"team.slug = {:teamSlug}",
		"-start",
		0,
		0,
		dbx.Params{"teamSlug": strings.ToLower(slug)},
	)
	if err != nil {
		return nil, err
	}

	if errs := t.dao.ExpandRecords(records, []string{"games_via_tournament"}, nil); len(errs) > 0 {
		return nil, fmt.Errorf("failed to expand: %s", errs)
	}

	tournaments := make([]*modelspb.TournamentWithGames, len(records))
	for i, r := range records {
		arr := r.ExpandedAll("games_via_tournament")
		slices.SortStableFunc(arr, func(a, b *models.Record) int {
			return a.GetDateTime("start_time").Time().Compare(b.GetDateTime("start_time").Time())
		})
		tournaments[i] = &modelspb.TournamentWithGames{
			Tournament: toTournament(r),
			Games:      toArr(arr, toGame),
		}
	}
	return tournaments, nil
}

func (t *Tournament) UpdateBySlug(teamSlug string, currentSlug string, name string, slug string, start types.DateTime, end types.DateTime, location string) (tournament *modelspb.Tournaments, err error) {

	if tournament, err = t.GetOneBySlug(teamSlug, currentSlug); err != nil {
		return nil, err
	}

	tournament.SetName(name)
	tournament.SetSlug(slug)
	tournament.SetStart(start)
	tournament.SetEnd(end)
	tournament.SetLocation(location)

	if err = t.dao.Save(tournament.Record); err != nil {
		return nil, err
	}
	return tournament, nil
}

func (t *Tournament) Create(teamId string, name string, slug string, start types.DateTime, end types.DateTime, location string) (*modelspb.Tournaments, error) {
	tournament := toTournament(models.NewRecord(t.collection))
	tournament.SetTeam(teamId)
	tournament.SetName(name)
	tournament.SetSlug(slug)
	tournament.SetStart(start)
	tournament.SetEnd(end)
	tournament.SetLocation(location)

	if err := t.dao.SaveRecord(tournament.Record); err != nil {
		return nil, err
	}
	return tournament, nil
}

func (t *Tournament) Delete(id string) error {
	if record, err := t.dao.FindRecordById(t.collection.Id, id); err != nil {
		return err
	} else if err = t.dao.DeleteRecord(record); err != nil {
		return err
	}
	return nil
}

func (t *Tournament) DeleteBySlug(teamSlug string, tournamentSlug string) error {
	if tournament, err := t.GetOneBySlug(teamSlug, tournamentSlug); err != nil {
		return err
	} else {
		return t.dao.DeleteRecord(tournament.Record)
	}
}

func toTournament(record *models.Record) *modelspb.Tournaments {
	return &modelspb.Tournaments{
		Record: record,
	}
}
