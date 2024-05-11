package repository

import (
	"strings"
	"ultigamecast/pbmodels"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
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

func (t *Tournament) Create(m *pbmodels.Tournaments) (error) {
	return t.dao.DB().Model(m).Exclude("Id").Insert();
}

func (t *Tournament) Update(m *pbmodels.Tournaments, attrs... string) (err error) {
	return t.dao.DB().Model(m).Exclude("Id", "Team").Update(attrs...)
}

func (t *Tournament) DeleteById(id string) error {
	_, err := t.dao.DB().Delete(t.collection.TableName(), dbx.HashExp{"id": id}).Execute()
	return err
}

func (t *Tournament) tournamentQuery() *dbx.SelectQuery {
	return t.dao.ModelQuery(&pbmodels.Tournaments{})
}
