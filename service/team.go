package service

import (
	"context"
	"ultigamecast/models"
)

type Team struct {
	q *models.Queries
}

func NewTeam(q *models.Queries) *Team {
	return &Team{
		q: q,
	}
}

func (t *Team) GetBySlug(slug string) (*models.Team, error) {
	if team, err := t.q.GetTeam(context.Background(), slug); err != nil {
		return nil, err
	} else {
		return &team, nil
	}
}

func CreateTeam() {
	
}
