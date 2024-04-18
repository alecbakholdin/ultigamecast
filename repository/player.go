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

func (p *Player) GetAllForTeamBySlug(slug string) ([]*modelspb.Players, error) {
	records, err := p.dao.FindRecordsByExpr("players", dbx.NewExp("team = (SELECT teams.id FROM teams WHERE LOWER(slug) = {:teamSlug} LIMIT 1)", dbx.Params{
		"teamSlug": strings.ToLower(slug),
	}))
	if err != nil {
		return nil, err
	}
	return toPlayers(records), nil
}

func toPlayers(records []*models.Record) []*modelspb.Players {
	players := make([]*modelspb.Players, len(records))
	for i, record := range records {
		players[i] = toPlayer(record)
	}
	return players
}

func toPlayer(record *models.Record) *modelspb.Players {
	return &modelspb.Players{
		Record: record,
	}
}
