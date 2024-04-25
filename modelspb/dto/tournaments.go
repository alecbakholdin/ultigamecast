package dto

import (
	"regexp"
	"strings"
	"ultigamecast/modelspb"
	"ultigamecast/validation"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/tools/types"
)

type Tournament struct {
	TeamSlug          string `param:"teamSlug"`
	TournamentSlug    string `param:"tournamentSlug"`
	TournamentSlugNew string
	Name              string `form:"name"`
	Start             string `form:"start"`
	StartDt           types.DateTime
	End               string `form:"end"`
	EndDt             types.DateTime
	Location          string `form:"location"`
}

func BindTournament(c echo.Context, payload *Tournament) (err error) {
	if err := c.Bind(payload); err != nil {
		return err
	}

	payload.TournamentSlugNew = convertToSlug(payload.Name)
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

func (t *Tournament) Load(m *modelspb.Tournaments) *Tournament {
	t.TournamentSlug = m.GetSlug()
	t.Name = m.GetName()
	t.Start = m.GetStart().Time().Format("2006-01-02")
	t.End = m.GetEnd().Time().Format("2006-01-02")
	t.Location = m.GetLocation()
	return t
}

func DtoFromTournament(t *modelspb.Tournaments) *Tournament {
	return &Tournament{
		TournamentSlug: t.GetSlug(),
		Name:           t.GetName(),
		Start:          t.GetStart().Time().Format("2006-01-02"),
		End:            t.GetEnd().Time().Format("2006-01-02"),
		Location:       t.GetLocation(),
	}
}

var whitespaceRegex = regexp.MustCompile(`[\s+]`)

func convertToSlug(s string) string {
	noWhitespace := whitespaceRegex.ReplaceAllString(s, "-")
	return strings.ToLower(noWhitespace)
}
