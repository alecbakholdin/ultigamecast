package pbmodels

import (
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/types"
)

type Tournaments struct {
	models.BaseModel

	Location      string         `db:"location" json:"location" form:"location"`
	Team          string         `db:"team" json:"team" form:"team"`
	Name          string         `db:"name" json:"name" form:"name"`
	Slug          string         `db:"slug" json:"slug" param:"tournamentsSlug"`
	Start         types.DateTime `db:"start" json:"start"`
	StartTimezone string         `json:"start_timezone" form:"start_timezone"`
	StartDatetime string         `json:"start_datetime" form:"start_datetime"`
	End           types.DateTime `db:"end" json:"end"`
	EndTimezone   string         `form:"end_timezone" json:"end_timezone"`
	EndDatetime   string         `form:"end_datetime" json:"end_datetime"`
}

func (m *Tournaments) TableName() string {
    return "tournaments"
}

