package pbmodels

import (
	"github.com/pocketbase/pocketbase/models"
)

type Players struct {
	models.BaseModel

	Order int    `db:"order" form:"order" json:"order"`
	Name  string `db:"name" form:"name" json:"name"`
	Team  string `db:"team" form:"team" json:"team"`
}


func (d *Players) CopyFrom(s *Players) *Players {
	d.Order = s.Order
	d.Name = s.Name
	d.Team = s.Team
	return d
}

func (m *Players) Copy() *Players {
	return (&Players{}).CopyFrom(m)
}

func (m *Players) TableName() string {
    return "players"
}
