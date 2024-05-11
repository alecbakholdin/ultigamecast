package repository

import (
	"ultigamecast/pbmodels"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
)

type Events struct {
	app        core.App
	dao        *daos.Dao
	tableName  string
	collection *models.Collection
}

func NewEvents(app core.App) *Events {
	return &Events{
		app:        app,
		dao:        app.Dao(),
		tableName:  "events",
		collection: mustGetCollection(app.Dao(), "events"),
	}
}

// Creates new event
func (e *Events) Create(event *pbmodels.Events) error {
	return e.dao.DB().Model(event).Exclude("Id").Insert()
}

// Updates an existing record. Make sure event has Id set or this will fail
// Updates only the fields specified in fields. If fields is empty, then updates
// everything (except for Id and Game, which should never change after creation)
func (e *Events) Update(event *pbmodels.Events, fields ...string) error {
	return e.dao.DB().Model(event).Exclude("Id", "Game").Update(fields...)
}

// Deletes an event by id
func (e *Events) Delete(eventId string) error {
	_, err := e.dao.DB().Delete(e.tableName, dbx.HashExp{"id": eventId}).Execute()
	return err
}
