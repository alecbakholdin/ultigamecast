package pbmodels

import (
	"github.com/pocketbase/pocketbase/models"
)

type Players struct {
	models.BaseModel

	Team  string `db:"team" form:"team" json:"team"`
	Name  string `db:"name" form:"name" json:"name"`
	Order int    `db:"order" form:"order" json:"order"`
}

func (m *Players) TableName() string {
    return "players"
}
