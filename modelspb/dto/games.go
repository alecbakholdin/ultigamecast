package dto

import (
	"fmt"
	"ultigamecast/modelspb"
	"ultigamecast/validation"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/tools/types"
)

type Games struct {
	TeamSlug          string `param:"teamSlug"`
	TournamentSlug    string `param:"tournamentSlug"`
	GameID            string `param:"gameId"`
	GameOpponent      string `form:"opponent"`
	GameTeamScore     int    `form:"team_score"`
	GameOpponentScore int    `form:"opponent_score"`
	GameHalfCap       int    `form:"half_cap"`
	GameSoftCap       int    `form:"soft_cap"`
	GameHardCap       int    `form:"hard_cap"`
	GameWindMph       int    `form:"wind_mph"`
	GameTempF         int    `form:"temp_f"`
	GameStartTime     string `form:"start_time"`
	GameStartTimeDt   types.DateTime
	GameStartTimezone string `form:"start_timezone"`
	GamesStatus       string `form:"status"`
}

func BindGameDto(c echo.Context, dto *Games) (err error) {
	if err := c.Bind(dto); err != nil {
		return err
	}

	if dto.GameStartTime != "" {
		if dto.GameStartTimeDt, err = types.ParseDateTime(dto.GameStartTime + ":00" + dto.GameStartTimezone); err != nil {
			c.Echo().Logger.Error(fmt.Errorf("error parsing %s: %s", dto.GameStartTime+dto.GameStartTimezone, err))
			validation.AddFieldErrorString(c, "start_time", "invalid date format")
		}
		fmt.Println(dto.GameStartTimeDt)
	}

	if dto.GameOpponent == "" {
		validation.AddFieldErrorString(c, "opponent", "cannot be empty")
	}

	return nil
}

func DtoFromGame(game *modelspb.Games) *Games {
	return &Games{
		GameID:            game.Record.GetId(),
		GameOpponent:      game.GetOpponent(),
		GameTeamScore:     game.GetTeamScore(),
		GameOpponentScore: game.GetOpponentScore(),
		GameHalfCap:       game.GetHalfCap(),
		GameSoftCap:       game.GetSoftCap(),
		GameHardCap:       game.GetHardCap(),
		GameWindMph:       game.GetWindMph(),
		GameTempF:         game.GetTempF(),
		GameStartTimeDt:   game.GetStartTime(),
		GameStartTime:     game.GetStartTime().String(),
		GamesStatus:       game.GetStatus(),
	}
}
