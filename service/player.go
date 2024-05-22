package service

import (
	"context"
	"database/sql"
	"ultigamecast/models"
)

type Player struct {
	q *models.Queries
	db *sql.DB
}

func NewPlayer(q *models.Queries, db *sql.DB) *Player {
	return &Player{
		q: q,
		db: db,
	}
}

func (p *Player) CreatePlayer(ctx context.Context, name string) {

}