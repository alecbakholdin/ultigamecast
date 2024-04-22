package models

import (
	"ultigamecast/validation"
	view "ultigamecast/view/team"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/tools/types"
)

type TournamentPayload struct {
	TournamentID   string
	TeamSlug       string `param:"teamSlug"`
	TournamentSlug string `param:"tournamentSlug"`
	Name           string `form:"name"`
	Start          string `form:"start"`
	StartDt        types.DateTime
	End            string `form:"end"`
	EndDt          types.DateTime
	Location       string `form:"location"`
}

func BindTournament(c echo.Context, payload *TournamentPayload) (err error) {
	if err := c.Bind(payload); err != nil {
		return err
	}

	payload.TournamentSlug = convertToSlug(payload.Name)
	if payload.Name == "" {
		validation.AddFieldErrorString(c, "name", "name cannot be empty")
	}
	if payload.StartDt, err = types.ParseDateTime(payload.Start); payload.Start != "" && err != nil {
		validation.AddFieldErrorString(c, "start", "invalid date format")
	}
	if payload.EndDt, err = types.ParseDateTime(payload.End); payload.End != "" && err != nil {
		validation.AddFieldErrorString(c, "end", "invalid date format")
	}
	if payload.StartDt.Time() != payload.EndDt.Time() && !payload.EndDt.Time().After(payload.StartDt.Time()) {
		validation.AddFieldErrorString(c, "end", "end date must be after start date")
	}

	return nil
}

func (p *TournamentPayload) ToData() view.TournamentData {
	return view.TournamentData{
		ID:       p.TournamentID,
		Name:     p.Name,
		Start:    p.Start,
		End:      p.End,
		Location: p.Location,
	}
}
