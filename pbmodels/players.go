package pbmodels

import (
	"github.com/pocketbase/pocketbase/models"
)

type Players struct {
	models.BaseModel

	Name  string `db:"name" form:"name" json:"name"`
	Team  string `db:"team" form:"team" json:"team"`
	Order int    `db:"order" form:"order" json:"order"`
}


func (d *Players) CopyFrom(s *Players) *Players {
	d.Name = s.Name
	d.Team = s.Team
	d.Order = s.Order
	return d
}
func (m *Players) TableName() string {
    return "players"
}
