package ctxvar

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"ultigamecast/internal/models"

	"github.com/a-h/templ"
)

type ContextVar string

const (
	HttpMethod      ContextVar = "method"
	Path            ContextVar = "path"
	ReqId           ContextVar = "RequestId"
	User            ContextVar = "User"
	Team            ContextVar = "Team"
	Player          ContextVar = "Player"
	Tournament      ContextVar = "Tournament"
	Game            ContextVar = "Game"
	Admin           ContextVar = "Admin"
	Event           ContextVar = "Event"
	TournamentDatum ContextVar = "TournamentDatum"
)

var LogMessageVars = []ContextVar{HttpMethod, Path}
var LogAttrVars = []ContextVar{ReqId, User}

func Url(ctx context.Context, segments ...any) string {
	urlParts := make([]string, 0)
	for _, cv := range segments {
		val := cv
		if ctxVal, ok := val.(ContextVar); ok {
			val = ctx.Value(ctxVal)
			if val == nil {
				panic(fmt.Sprintf("missing context value for key %s", ctxVal))
			}
		}

		switch v := val.(type) {
		case string:
			if v = strings.TrimSpace(v); len(v) > 0 {
				urlParts = append(urlParts, v)
			}
		case *models.Team:
			urlParts = append(urlParts, "teams", v.Slug)
		case *models.TournamentSummary:
			urlParts = append(urlParts, "schedule", v.Slug)
		case *models.Tournament:
			urlParts = append(urlParts, "schedule", v.Slug)
		case *models.Player:
			urlParts = append(urlParts, "players", v.Slug)
		case *models.Game:
			urlParts = append(urlParts, "schedule", v.Slug)
		case *models.Event:
			urlParts = append(urlParts, "events", v.ID)
		case *models.TournamentDatum:
			urlParts = append(urlParts, "data", strconv.FormatInt(v.ID, 10))
		default:
			panic(fmt.Sprintf("unexpected type %T", v))
		}
	}
	return "/" + strings.Join(urlParts, "/")
}

func SafeUrl(ctx context.Context, segments ...any) templ.SafeURL {
	return templ.SafeURL(Url(ctx, segments...))
}

func IsAdmin(ctx context.Context) bool {
	if val, ok := ctx.Value(Admin).(bool); ok {
		return val
	} else {
		return false
	}
}

func IsAuthenticated(ctx context.Context) bool {
	return GetUser(ctx) != nil
}

func GetUser(ctx context.Context) *models.User {
	return getModel[models.User](ctx, User)
}

func GetTeam(ctx context.Context) *models.Team {
	return getModel[models.Team](ctx, Team)
}

func GetPlayer(ctx context.Context) *models.Player {
	return getModel[models.Player](ctx, Player)
}

func GetTournament(ctx context.Context) *models.Tournament {
	return getModel[models.Tournament](ctx, Tournament)
}

func GetGame(ctx context.Context) *models.Game {
	return getModel[models.Game](ctx, Game)
}

func GetEvent(ctx context.Context) *models.Event {
	return getModel[models.Event](ctx, Event)
}

func GetTournamentDatum(ctx context.Context) *models.TournamentDatum {
	return getModel[models.TournamentDatum](ctx, TournamentDatum)
}

func getModel[T interface{}](ctx context.Context, key ContextVar) *T {
	if m, ok := ctx.Value(key).(*T); ok {
		return m
	}
	return nil
}

func GetValue(ctx context.Context, key ContextVar) string {
	switch key {
	case User:
		if user := GetUser(ctx); user != nil {
			return fmt.Sprintf("[%d] %s", user.ID, user.Email)
		}
	case Team:
		if team := GetTeam(ctx); team != nil {
			return fmt.Sprintf("[%d] %s", team.ID, team.Name)
		}
	case Player:
		if player := GetPlayer(ctx); player != nil {
			return fmt.Sprintf("[%d] %s", player.ID, player.Name)
		}
	case Tournament:
		if tournament := GetTournament(ctx); tournament != nil {
			return fmt.Sprintf("[%d] %s", tournament.ID, tournament.Name)
		}
	case Game:
		if game := GetGame(ctx); game != nil {
			return fmt.Sprintf("[%d] %s", game.ID, game.Opponent)
		}
	case Event:
		if event := GetEvent(ctx); event != nil {
			return fmt.Sprintf("[%s] %s", event.ID, event.Type)
		}
	case TournamentDatum:
		if datum := GetTournamentDatum(ctx); datum != nil {
			return fmt.Sprintf("[%d] %s", datum.ID, datum.Title)
		}

	default:
		if val := ctx.Value(key); val != nil {
			return fmt.Sprintf("%#v", val)
		}
	}
	return ""
}
