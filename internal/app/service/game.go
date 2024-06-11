package service

import (
	"context"
	"database/sql"
	"ultigamecast/internal/models"
)

type Game struct {
	q  *models.Queries
	db *sql.DB
}

func NewGame(q *models.Queries, db *sql.DB) *Game {
	return &Game{
		q:  q,
		db: db,
	}
}

func (g *Game) GetGame(ctx context.Context, slug string) (*models.Game, error) {
	return nil, nil
}

func (g *Game) CreateGame(ctx context.Context) {

}
