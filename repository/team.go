package repository

import (
	"fmt"
	"path"
	"ultigamecast/modelspb"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
	"gocloud.dev/blob"
)

type Team struct {
	app        core.App
	dao        *daos.Dao
	collection *models.Collection
}

func NewTeam(app core.App) *Team {
	return &Team{
		app:        app,
		dao:        app.Dao(),
		collection: mustGetCollection(app.Dao(), "teams"),
	}
}

func (t *Team) GetOneBySlug(slug string) (*modelspb.Teams, error) {
	if record, err := t.dao.FindFirstRecordByData(t.collection.Id, "slug", slug); err != nil {
		return nil, err
	} else {
		return toTeam(record), nil
	}
}

func (t *Team) GetLogo(slug string) (*blob.Reader, error) {
	if team, err := t.GetOneBySlug(slug); err != nil {
		return nil, err
	} else if logo := team.GetLogo(); logo == "" {
		return nil, fmt.Errorf("logo doesnt exist")
	} else if filesystem, err := t.app.NewFilesystem(); err != nil {
		return nil, err
	} else if reader, err := filesystem.GetFile(path.Join(team.Record.BaseFilesPath(), logo)); err != nil {
		return nil, err
	} else {
		return reader, nil
	}
}

func toTeam(record *models.Record) *modelspb.Teams {
	return &modelspb.Teams{
		Record: record,
	}
}
