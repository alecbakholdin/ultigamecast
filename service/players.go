package service

import (
	"fmt"
	"ultigamecast/pbmodels"
	"ultigamecast/repository"
)

type Players struct {
	PlayerRepo *repository.Player
	TeamRepo   *repository.Team
}

func NewPlayers(p *repository.Player, t *repository.Team) *Players {
	return &Players{
		PlayerRepo: p,
		TeamRepo:   t,
	}
}

func (p *Players) GetAllBySlug(teamSlug string) ([]*pbmodels.Players, error) {
	if players, err := p.PlayerRepo.GetAllByTeamSlug(teamSlug); err != nil {
		return nil, fmt.Errorf("error fetching players for team [%s]: %s", teamSlug, err)
	} else {
		return players, nil
	}
}

func (p *Players) GetOneById(id string) (*pbmodels.Players, error) {
	if player, err := p.PlayerRepo.GetOneById(id); err != nil {
		return nil, err
	} else {
		return player, nil
	}
}

func (p *Players) Create(teamSlug string, player *pbmodels.Players) error {
	if teamSlug == "" {
		return fmt.Errorf("error creating player: teamSlug cannot be empty")
	}

	var (
		team *pbmodels.Teams
		err  error
	)
	if team, err = p.TeamRepo.FindOneBySlug(teamSlug); err != nil {
		return fmt.Errorf("error finding team %s while creating player: %s", teamSlug, err)
	}

	player.Team = team.Id
	if err = p.PlayerRepo.Create(player); err != nil {
		return fmt.Errorf("error creating player for team %s: %s", teamSlug, err)
	}
	return nil
}

func (p *Players) Update(playerId string, player *pbmodels.Players) error {
	if playerId == "" {
		return fmt.Errorf("playerId cannot be empty")
	}
	player.Id = playerId
	return p.PlayerRepo.Update(player)
}

func (p *Players) UpdateOrder(playerIdOrder []string) error {
	err := p.PlayerRepo.UpdateOrder(playerIdOrder)
	if err != nil {
		return fmt.Errorf("error updating player order: %s", err)
	}
	return nil
}

func (p *Players) Delete(playerId string) error {
	if err := p.PlayerRepo.Delete(playerId); err != nil {
		return fmt.Errorf("error deleting player [%s]: %s", playerId, err)
	}
	return nil
}