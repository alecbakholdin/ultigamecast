package repository

import (
	"strings"
	"ultigamecast/modelspb"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/forms"
	"github.com/pocketbase/pocketbase/models"
)

type Player struct {
	app                   core.App
	dao                   *daos.Dao
	collection            *models.Collection
	tournamentSummaryView *models.Collection
	teamSummaryView       *models.Collection
	gameSummaryView       *models.Collection
}

func NewPlayer(app core.App) *Player {
	dao := app.Dao()
	return &Player{
		app:                   app,
		dao:                   dao,
		collection:            mustGetCollection(dao, "players"),
		tournamentSummaryView: mustGetCollection(dao, "player_tournament_summary"),
		teamSummaryView:       mustGetCollection(dao, "player_team_summary"),
		gameSummaryView:       mustGetCollection(dao, "player_game_summary"),
	}
}

func (p *Player) GetOneById(id string) (*modelspb.Players, error) {
	if record, err := p.dao.FindRecordById(p.collection.Id, id); err != nil {
		return nil, err
	} else {
		return toPlayer(record), nil
	}
}

func (p *Player) GetAllByTeamSlug(slug string) ([]*modelspb.Players, error) {
	records, err := p.dao.FindRecordsByFilter(p.collection.Id, "team.slug = {:teamSlug}", "+order", 0, 0, dbx.Params{"teamSlug": strings.ToLower(slug)})
	if err != nil {
		return nil, err
	}
	return toArr(records, toPlayer), nil
}

func (p *Player) Create(team *modelspb.Teams, name string, order int) (*modelspb.Players, error) {
	player := toPlayer(models.NewRecord(p.collection))
	player.SetTeam(team.Record.Id)
	player.SetName(name)
	player.SetOrder(order)

	form := forms.NewRecordUpsert(p.app, player.Record)

	if err := form.Submit(); err != nil {
		return nil, err
	}
	return player, nil
}

func (p *Player) Update(id string, name string) (*modelspb.Players, error) {
	var player *modelspb.Players
	if record, err := p.dao.FindRecordById(p.collection.Id, id); err != nil {
		return nil, err
	} else {
		player = toPlayer(record)
	}

	player.SetName(name)

	if err := p.dao.SaveRecord(player.Record); err != nil {
		return nil, err
	}
	return player, nil
}

func (p *Player) Delete(id string) error {
	if record, err := p.dao.FindRecordById(p.collection.Id, id); err != nil {
		return err
	} else {
		return p.dao.DeleteRecord(record)
	}
}

func (p *Player) UpdateOrder(players []*modelspb.Players, playerIds []string) ([]*modelspb.Players, error) {
	sortedPlayers := make([]*modelspb.Players, len(players))
	for order, id := range playerIds {
		for _, player := range players {
			if player.Record.GetId() == id {
				sortedPlayers[order] = player
				break
			}
		}
	}
	err := p.dao.RunInTransaction(func(txDao *daos.Dao) error {
		for order, player := range sortedPlayers {
			player.SetOrder(-(order + 1))
			if err := txDao.SaveRecord(player.Record); err != nil {
				return err
			}
		}
		for order, player := range sortedPlayers {
			player.SetOrder(order)
			if err := txDao.SaveRecord(player.Record); err != nil {
				return err
			}
		}
		return nil
	})
	return sortedPlayers, err
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
