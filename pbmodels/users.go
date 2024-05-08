package pbmodels

import (
	"github.com/pocketbase/pocketbase/models"
	"path"
)

type Users struct {
	models.BaseModel

	Name   string `db:"name" form:"name" json:"name"`
	Avatar string `db:"avatar" form:"avatar" json:"avatar"`
}

func (m *Users) TableName() string {
    return "users"
}

func (m *Users) GetAvatarPath() string {
	if m.Avatar == "" || m.GetId() == "" {
		return ""
	} else {
		return path.Join("_pb_users_auth_", m.GetId(), m.Avatar)
	}
}
