package repository

import (
	"strings"
	"ultigamecast/modelspb"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
)

type Player struct {
	dao *daos.Dao
}

func NewPlayer(dao *daos.Dao) *Player {
	return &Player{
		dao: dao,
	}
}

func (p *Player) GetAllByTeamSlug(slug string) ([]*modelspb.Players, error) {
	records, err := p.dao.FindRecordsByFilter("players", "team.slug = {:teamSlug}", "+order", 0, 0, dbx.Params{"teamSlug": strings.ToLower(slug)})
	if err != nil {
		return nil, err
	}
	return toArr(records, toPlayer), nil
}

func toPlayer(record *models.Record) *modelspb.Players {
	return &modelspb.Players{
		Record: record,
	}
}
