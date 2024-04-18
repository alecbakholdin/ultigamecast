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

func (p *Player) GetAllByTeamSlug(teamSlug string) ([]*modelspb.Players, error) {
	records, err := p.dao.FindRecordsByExpr("players", dbx.NewExp("team = (SELECT teams.id FROM teams WHERE LOWER(slug) = {:teamSlug} LIMIT 1)", dbx.Params{
		"teamSlug": strings.ToLower(teamSlug),
	}))
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
