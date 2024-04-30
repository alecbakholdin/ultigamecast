package repository

import (
	"strings"
	"ultigamecast/modelspb"
	"ultigamecast/pbmodels"

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

func (t *Tournament) GetOneBySlug(teamSlug string, tournamentSlug string) (*pbmodels.Tournaments, error) {
	tournament := &pbmodels.Tournaments{}
	err := t.tournamentQuery().InnerJoin("teams", dbx.NewExp("teams.id = tournaments.team")).Where(dbx.HashExp{
		"tournaments.slug": strings.ToLower(tournamentSlug),
		"teams.slug":       strings.ToLower(teamSlug),
	}).One(tournament)
	if err != nil {
		return nil, err
	}
	return tournament, err
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

func (t *Tournament) GetAllByTeamSlug(slug string) ([]*pbmodels.Tournaments, error) {
	q := t.tournamentQuery()
	q.InnerJoin("teams", dbx.NewExp("teams.id = tournaments.team"))
	q.Where(dbx.HashExp{"teams.slug": slug})
	q.OrderBy("start DESC")

	tournaments := make([]*pbmodels.Tournaments, 0)
	if err := q.All(&tournaments); err != nil {
		return nil, err
	}
	return tournaments, nil
}

func (t *Tournament) UpdateBySlug(teamSlug string, currentSlug string, name string, slug string, start types.DateTime, end types.DateTime, location string) (*pbmodels.Tournaments, error) {
	currentModel, err := t.GetOneBySlug(teamSlug, currentSlug)
	if err != nil {
		return nil, err
	}

	currentModel.Name = name
	currentModel.Slug = slug
	currentModel.Start = start
	currentModel.End = end
	currentModel.Location = location

	if err := t.dao.DB().Model(currentModel).Update(); err != nil {
		return nil, err
	} else {
		return currentModel, nil
	}
}

func (t *Tournament) Create(teamId string, name string, slug string, start types.DateTime, end types.DateTime, location string) (*pbmodels.Tournaments, error) {
	tournament := &pbmodels.Tournaments{
		Team:     teamId,
		Name:     name,
		Slug:     slug,
		Start:    start,
		End:      end,
		Location: location,
	}
	if err := t.dao.DB().Model(tournament).Insert(); err != nil {
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
	if _, err := t.GetOneBySlug(teamSlug, tournamentSlug); err != nil {
		return err
	} else {
		panic("Tournament is not deletable rn") //TODO
	}
}

func toTournament(record *models.Record) *modelspb.Tournaments {
	return &modelspb.Tournaments{
		Record: record,
	}
}

func (t *Tournament) tournamentQuery() *dbx.SelectQuery {
	return t.dao.ModelQuery(&pbmodels.Tournaments{})
}
