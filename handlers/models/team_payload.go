
package models

import (
	view "ultigamecast/view/team"

	"github.com/labstack/echo/v5"
)

type TeamPayload struct {
	TeamSlug       string `param:"teamSlug"`
	Name           string `form:"name"`

}

func BindTeam(c echo.Context, payload *TeamPayload) (err error) {
	if err := c.Bind(payload); err != nil {
		return err
	}

	return nil
}
func (p *TeamPayload) ToData() view.TeamData {
	return view.TeamData{
		TeamSlug: p.TeamSlug,
		TeamName: p.Name,
	}
}
