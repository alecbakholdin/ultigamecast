package service

import "ultigamecast/models"

type Team struct {
	q *models.Queries
}

func NewTeam(q *models.Queries) *Team {
	return &Team{
		q: q,
	}
}

func CreateTeam() {
	
}