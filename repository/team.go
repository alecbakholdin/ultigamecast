package repository

import (
	"ultigamecast/modelspb"

	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
)

const (
	collectionName = "teams"
	slugField      = "slug"
)

type Team struct {
	dao *daos.Dao
}

func NewTeam(dao *daos.Dao) *Team {
	return &Team{
		dao: dao,
	}
}

func (t *Team) GetOneBySlug(slug string) (*modelspb.Teams, error) {
	if record, err := t.dao.FindFirstRecordByData(collectionName, slugField, slug); err != nil {
		return nil, err
	} else {
		return toTeam(record), nil
	}
}

func toTeam(record *models.Record) *modelspb.Teams {
	return &modelspb.Teams{
		Record: record,
	}
}
