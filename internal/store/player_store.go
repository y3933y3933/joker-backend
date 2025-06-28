package store

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/y3933y3933/joker/internal/db/sqlc"
)

type Player struct {
	ID       int64  `json:"id"`
	Nickname string `json:"nickname"`
	IsHost   bool   `json:"isHost"`
	GameID   int64  `json:"gameID"`
}

type PostgresPlayerStore struct {
	queries *sqlc.Queries
}

func NewPostgresPlayerStore(queries *sqlc.Queries) *PostgresPlayerStore {
	return &PostgresPlayerStore{queries: queries}
}

type PlayerStore interface {
	Create(ctx context.Context, player *Player) (*Player, error)
	CountPlayerInGame(ctx context.Context, gameID int64) (int64, error)
	FindPlayersByGameID(ctx context.Context, gameID int64) ([]*Player, error)
}

func (pg *PostgresPlayerStore) Create(ctx context.Context, player *Player) (*Player, error) {
	args := sqlc.CreatePlayerParams{
		GameID:   player.GameID,
		Nickname: player.Nickname,
		IsHost: pgtype.Bool{
			Bool:  player.IsHost,
			Valid: true,
		},
	}

	row, err := pg.queries.CreatePlayer(ctx, args)

	if err != nil {
		return nil, err
	}

	player.ID = row.ID
	return player, nil

}

func (pg *PostgresPlayerStore) CountPlayerInGame(ctx context.Context, gameID int64) (int64, error) {
	count, err := pg.queries.CountPlayersInGame(ctx, gameID)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (pg *PostgresPlayerStore) FindPlayersByGameID(ctx context.Context, gameID int64) ([]*Player, error) {
	dbPlayers, err := pg.queries.FindPlayersByGameID(ctx, gameID)
	if err != nil {
		return nil, err
	}

	players := make([]*Player, 0, len(dbPlayers))
	for _, p := range dbPlayers {
		players = append(players, &Player{
			ID:       p.ID,
			Nickname: p.Nickname,
			IsHost:   p.IsHost.Bool,
			GameID:   p.GameID,
		})
	}
	return players, nil
}
