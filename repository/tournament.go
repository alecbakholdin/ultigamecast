package repository

import (
	"strings"
	"ultigamecast/modelspb"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/types"
)

type Tournament struct {
	dao        *daos.Dao
	collection *models.Collection
}

func NewTournament(dao *daos.Dao) *Tournament {
	collection, err := dao.FindCollectionByNameOrId("tournaments")
	if err != nil {
		panic(err)
	}
	return &Tournament{
		dao:        dao,
		collection: collection,
	}
}

func (t *Tournament) GetOneBySlug(teamSlug string, tournamentSlug string) (*modelspb.Tournaments, error) {
	record, err := t.dao.FindFirstRecordByFilter(
		t.collection.Name,
		"team.slug = {:teamSlug} && slug = {:tournamentSlug}",
		dbx.Params{
			"teamSlug": teamSlug,
			"tournamentSlug": tournamentSlug,
		},
	)
	if err != nil {
		return nil, err
	}
	return toTournament(record), err
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

func (t *Tournament) Create(team *modelspb.Teams, name string, slug string, start types.DateTime, end types.DateTime, location string) (*modelspb.Tournaments, error) {
	tournament := toTournament(models.NewRecord(t.collection))
	tournament.SetTeam(team.Record.GetId())
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

func toTournament(record *models.Record) *modelspb.Tournaments {
	return &modelspb.Tournaments{
		Record: record,
	}
}
