package repository

import (
	"strings"
	"ultigamecast/modelspb"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
)

type Tournament struct {
	dao *daos.Dao
}

func NewTournament(dao *daos.Dao) *Tournament {
	return &Tournament{
		dao: dao,
	}
}

func (p *Tournament) GetAllByTeamSlug(slug string) ([]*modelspb.Tournaments, error) {
	records, err := p.dao.FindRecordsByFilter(
		"tournaments",
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

func toTournament(record *models.Record) *modelspb.Tournaments {
	return &modelspb.Tournaments{
		Record: record,
	}
}
