package repository

import (
	"strings"
	"ultigamecast/modelspb"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
)

type Player struct {
	dao                   *daos.Dao
	collection            *models.Collection
	tournamentSummaryView *models.Collection
	teamSummaryView       *models.Collection
	gameSummaryView       *models.Collection
}

func NewPlayer(dao *daos.Dao) *Player {
	collection, err := dao.FindCollectionByNameOrId("players")
	if err != nil {
		panic(err)
	}
	tournamentSummaryView, err := dao.FindCollectionByNameOrId("player_tournament_summary")
	if err != nil {
		panic(err)
	}
	teamSummaryView, err := dao.FindCollectionByNameOrId("player_team_summary")
	if err != nil {
		panic(err)
	}
	gameSummaryView, err := dao.FindCollectionByNameOrId("player_game_summary")
	if err != nil {
		panic(err)
	}

	return &Player{
		dao:                   dao,
		collection:            collection,
		tournamentSummaryView: tournamentSummaryView,
		teamSummaryView:       teamSummaryView,
		gameSummaryView:       gameSummaryView,
	}
}

func (p *Player) GetAllByTeamSlug(slug string) ([]*modelspb.Players, error) {
	records, err := p.dao.FindRecordsByFilter(p.collection.Id, "team.slug = {:teamSlug}", "+order", 0, 0, dbx.Params{"teamSlug": strings.ToLower(slug)})
	if err != nil {
		return nil, err
	}
	return toArr(records, toPlayer), nil
}

type PlayerSummaryType string

const (
	PlayerSummaryTypeTeam       PlayerSummaryType = "team"
	PlayerSummaryTypeGame       PlayerSummaryType = "game"
	PlayerSummaryTypeTournament PlayerSummaryType = "tournament"
)

func (p *Player) GetPlayerTeamSummariesByTeamSlug(teamSlug string, orderByField string, direction SortDirection) ([]modelspb.PlayerSummary, error) {
	records, err := p.dao.FindRecordsByFilter(
		p.teamSummaryView.Id,
		"team_slug = {:teamSlug}",
		GetSortString(direction, orderByField),
		0,
		0,
		dbx.Params{"teamSlug": strings.ToLower(teamSlug)},
	)
	if err != nil && !IsNotFound(err) {
		return nil, err
	} else if IsNotFound(err) {
		return []modelspb.PlayerSummary{}, nil
	}
	return toArr(records, toPlayerSummary), nil
}

func (p *Player) GetPlayerTournamentSummariesByTeamAndTournamentSlugs(teamSlug string, tournamentSlug string, orderByField string, direction SortDirection) ([]modelspb.PlayerSummary, error) {
	records, err := p.dao.FindRecordsByFilter(
		p.tournamentSummaryView.Id,
		"team_slug = {:teamSlug} && tournament_slug = {:tournamentSlug}",
		GetSortString(direction, orderByField),
		0,
		0,
		dbx.Params{"teamSlug": strings.ToLower(teamSlug), "tournamentSlug": strings.ToLower(tournamentSlug)},
	)
	if err != nil && !IsNotFound(err) {
		return nil, err
	} else if IsNotFound(err) {
		return []modelspb.PlayerSummary{}, nil
	}
	return toArr(records, toTournamentSummary), nil
}

func (p *Player) GetPlayerGameSummariesByGameId(gameId string, orderByField string, direction SortDirection) ([]modelspb.PlayerSummary, error) {
	records, err := p.dao.FindRecordsByFilter(
		p.tournamentSummaryView.Id,
		"game_id = {:gameId}",
		GetSortString(direction, orderByField),
		0,
		0,
		dbx.Params{"gameId": gameId},
	)
	if err != nil && !IsNotFound(err) {
		return nil, err
	} else if IsNotFound(err) {
		return []modelspb.PlayerSummary{}, nil
	}
	return toArr(records, toGameSummary), nil
}

func toPlayer(record *models.Record) *modelspb.Players {
	return &modelspb.Players{
		Record: record,
	}
}

func toPlayerSummary(record *models.Record) modelspb.PlayerSummary {
	return &modelspb.PlayerTeamSummary{
		Record: record,
	}
}

func toTournamentSummary(record *models.Record) modelspb.PlayerSummary {
	return &modelspb.PlayerTournamentSummary{
		Record: record,
	}
}

func toGameSummary(record *models.Record) modelspb.PlayerSummary {
	return &modelspb.PlayerGameSummary{
		Record: record,
	}
}
