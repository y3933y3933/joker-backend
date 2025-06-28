package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/y3933y3933/joker/internal/db/sqlc"
)

type Game struct {
	ID     int64  `json:"id"`
	Code   string `json:"code"`
	Status string `json:"status"`
}

const (
	GameStatusWaiting = "waiting"
	GameStatusPlaying = "playing"
	GameStatusEnded   = "ended"
)

type PostgresGameStore struct {
	queries *sqlc.Queries
}

func NewPostgresGameStore(queries *sqlc.Queries) *PostgresGameStore {
	return &PostgresGameStore{queries: queries}
}

type GameStore interface {
	Create(*Game) (*Game, error)
	GameCodeExists(code string) (bool, error)
}

func (pg *PostgresGameStore) Create(game *Game) (*Game, error) {
	ctx := context.Background()
	args := sqlc.CreateGameParams{
		Code:   game.Code,
		Status: game.Status,
	}
	row, err := pg.queries.CreateGame(ctx, args)
	if err != nil {
		return nil, err
	}

	return &Game{ID: row.ID, Code: game.Code, Status: game.Status}, nil
}

func (pg *PostgresGameStore) GameCodeExists(code string) (bool, error) {
	ctx := context.Background()
	_, err := pg.queries.GetGameByCode(ctx, code)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return true, nil

}
