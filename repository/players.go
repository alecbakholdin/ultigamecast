package repository

import (
	"fmt"
	"strings"
	"ultigamecast/pbmodels"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
)

type Player struct {
	app                   core.App
	dao                   *daos.Dao
	collection            *models.Collection
}

func NewPlayer(app core.App) *Player {
	dao := app.Dao()
	return &Player{
		app:                   app,
		dao:                   dao,
		collection:            mustGetCollection(dao, "players"),
	}
}

func (p *Player) GetOneById(id string) (*pbmodels.Players, error) {
	player := pbmodels.Players{}
	err := p.dao.ModelQuery(&pbmodels.Players{}).Where(dbx.HashExp{"id": id}).One(player)
	if err != nil {
		return nil, err
	}
	return &player, err
}

func (p *Player) GetAllByTeamSlug(slug string) ([]*pbmodels.Players, error) {
	q := p.dao.ModelQuery(&pbmodels.Players{})
	q.InnerJoin("teams", dbx.NewExp("teams.id = players.team"))
	q.Where(dbx.HashExp{"teams.slug": strings.ToLower(slug)})
	q.OrderBy("order ASC")

	players := make([]*pbmodels.Players, 0)
	if err := q.All(&players); err != nil {
		return nil, err
	}
	return players, nil
}

func (p *Player) Create(player *pbmodels.Players) error {
	return p.dao.DB().Model(player).Exclude("Id").Insert()
}

func (p *Player) Update(player *pbmodels.Players) error {
	return p.dao.DB().Model(player).Exclude("Team").Update()
}

func (p *Player) Delete(id string) error {
	_, err := p.dao.DB().Delete("players", dbx.HashExp{"id": id}).Execute()
	return err
}

func (p *Player) UpdateFieldBulk(playerIds []string, values []any, field string) error {
	if len(playerIds) != len(values) {
		return fmt.Errorf("playerIds and values must be same length but found playerIds [%d] values [%d]", len(playerIds), len(values))
	}

	params := dbx.Params{}
	playerIdList := make([]string, 2)
	caseStatement := "CASE id "
	for i, playerId := range playerIds {
		idStr := fmt.Sprintf("id_%d", i)
		valueStr := fmt.Sprintf("value_%d", i)
		caseStatement += fmt.Sprintf(" WHEN {:%s} THEN {:%s} ", idStr, valueStr)

		params[idStr] = playerId
		params[valueStr] = values[i]

		playerIdList[i] = fmt.Sprintf("{:%s}", idStr)
		if i == 1 {
			break
		}
	}
	caseStatement += " END "

	q := p.dao.DB().NewQuery(fmt.Sprintf("UPDATE players SET [%s] = %s WHERE id IN (%s)", field, caseStatement, strings.Join(playerIdList, ",")))
	q.Bind(params)
	_, err := q.Execute()
	return err
}

func (p *Player) UpdateOrder(playerOrder []string) error {
	return p.dao.RunInTransaction(func(txDao *daos.Dao) error {
		db := txDao.DB()
		for i, playerId := range playerOrder {
			q := db.Update("players", dbx.Params{"order": -(i + 1)}, dbx.HashExp{"id": playerId})
			if _, err := q.Execute(); err != nil {
				return fmt.Errorf("error updating playerId %s to order %d: %s", playerId, i, err)
			}
		}
		for i, playerId := range playerOrder {
			q := db.Update("players", dbx.Params{"order": i}, dbx.HashExp{"id": playerId})
			if _, err := q.Execute(); err != nil {
				return fmt.Errorf("error updating playerId %s to order %d: %s", playerId, i, err)
			}
		}
		return nil
	})
}
