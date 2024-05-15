package pbmodels

import (
	"github.com/pocketbase/pocketbase/models"
	"path"
)

type Teams struct {
	models.BaseModel

	Organization string `db:"organization" form:"organization" json:"organization"`
	Logo         string `db:"logo" form:"logo" json:"logo"`
	Name         string `db:"name" form:"name" json:"name"`
	Slug         string `db:"slug" json:"slug" param:"teamsSlug"`
}


func (d *Teams) CopyFrom(s *Teams) *Teams {
	d.Organization = s.Organization
	d.Logo = s.Logo
	d.Name = s.Name
	d.Slug = s.Slug
	return d
}

func (m *Teams) Copy() *Teams {
	return (&Teams{}).CopyFrom(m)
}

func (m *Teams) TableName() string {
    return "teams"
}

func (m *Teams) GetLogoPath() string {
	if m.Logo == "" || m.GetId() == "" {
		return ""
	} else {
		return path.Join("4edlnrnrqy9uk5q", m.GetId(), m.Logo)
	}
}
