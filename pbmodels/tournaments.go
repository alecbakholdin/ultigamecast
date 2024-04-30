package pbmodels

import (
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/types"
)

type Tournaments struct {
	models.BaseModel

	Slug string `db:"slug" json:"slug" param:"tournamentsSlug"`
	Start types.DateTime `db:"start" json:"start"`
	StartTimezone string `form:"start_timezone" json:"start_timezone"`
	StartDatetime string `form:"start_datetime" json:"start_datetime"`
	End types.DateTime `db:"end" json:"end"`
	EndTimezone string `form:"end_timezone" json:"end_timezone"`
	EndDatetime string `json:"end_datetime" form:"end_datetime"`
	Location string `db:"location" json:"location" form:"location"`
	Team string `db:"team" json:"team" form:"team"`
	Name string `db:"name" json:"name" form:"name"`
}

func (m *Tournaments) TableName() string {
    return "tournaments"
}

