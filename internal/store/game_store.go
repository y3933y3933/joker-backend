package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/y3933y3933/joker/internal/db/sqlc"
	"github.com/y3933y3933/joker/internal/utils/errx"
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
	Create(context.Context, *Game) (*Game, error)
	GameCodeExists(ctx context.Context, code string) (bool, error)
	GetGameByCode(ctx context.Context, code string) (*Game, error)
	UpdateStatus(ctx context.Context, gameID int64, status string) error
	EndGame(ctx context.Context, code string) error
}

func (pg *PostgresGameStore) Create(ctx context.Context, game *Game) (*Game, error) {
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

func (pg *PostgresGameStore) GameCodeExists(ctx context.Context, code string) (bool, error) {
	_, err := pg.queries.GetGameByCode(ctx, code)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return true, nil

}

func (pg *PostgresGameStore) GetGameByCode(ctx context.Context, code string) (*Game, error) {
	game, err := pg.queries.GetGameByCode(ctx, code)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errx.ErrGameNotFound
		}
		return nil, err
	}

	return &Game{
		ID:     game.ID,
		Code:   game.Code,
		Status: game.Status,
	}, nil
}

func (pg *PostgresGameStore) UpdateStatus(ctx context.Context, gameID int64, status string) error {
	return pg.queries.UpdateGameStatus(ctx, sqlc.UpdateGameStatusParams{
		ID:     gameID,
		Status: string(status),
	})
}

func (pg *PostgresGameStore) EndGame(ctx context.Context, code string) error {
	return pg.queries.EndGame(ctx, code)
}
