package pbmodels

import (
	"github.com/pocketbase/pocketbase/models"
)

type Players struct {
	models.BaseModel

	Team string `db:"team" json:"team" form:"team"`
	Name string `db:"name" json:"name" form:"name"`
	Order int `db:"order" json:"order" form:"order"`
}

func (m *Players) TableName() string {
    return "players"
}

