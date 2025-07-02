package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/y3933y3933/joker/internal/db/sqlc"
	"github.com/y3933y3933/joker/internal/utils/errx"
)

type Player struct {
	ID       int64  `json:"id"`
	Nickname string `json:"nickname"`
	IsHost   bool   `json:"isHost"`
	GameID   int64  `json:"gameID"`
	Status   string `json:"status"`
}

const (
	PlayerStatusOnline  = "online"
	PlayerStatusOffline = "offline"
)

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
	DeleteByID(ctx context.Context, id int64) error
	FindByID(ctx context.Context, id int64) (*Player, error)
	UpdateHost(ctx context.Context, id int64, isHost bool) error
	FindByNickname(ctx context.Context, gameID int64, nickname string) (*Player, error)
	GetPlayerCountByGameCode(ctx context.Context, gameCode string) (int64, error)
}

func (pg *PostgresPlayerStore) Create(ctx context.Context, player *Player) (*Player, error) {
	args := sqlc.CreatePlayerParams{
		GameID:   player.GameID,
		Nickname: player.Nickname,
		IsHost:   toPgBool(&player.IsHost),
	}

	row, err := pg.queries.CreatePlayer(ctx, args)

	if err != nil {
		return nil, err
	}

	return &Player{
		ID:       row.ID,
		Nickname: row.Nickname,
		IsHost:   fromPgBool(row.IsHost),
		GameID:   row.GameID,
		Status:   row.Status,
	}, nil

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
			Status:   p.Status,
		})
	}
	return players, nil
}

func (pg *PostgresPlayerStore) DeleteByID(ctx context.Context, id int64) error {
	return pg.queries.DeletePlayerByID(ctx, id)
}

func (pg *PostgresPlayerStore) FindByID(ctx context.Context, id int64) (*Player, error) {
	res, err := pg.queries.FindPlayerByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errx.ErrPlayerNotFound
		}
		return nil, err
	}

	return &Player{
		ID:       res.ID,
		Nickname: res.Nickname,
		IsHost:   fromPgBool(res.IsHost),
		GameID:   res.GameID,
		Status:   res.Status,
	}, nil
}

func (pg *PostgresPlayerStore) UpdateHost(ctx context.Context, id int64, isHost bool) error {
	return pg.queries.UpdateHost(ctx, sqlc.UpdateHostParams{
		ID:     id,
		IsHost: toPgBool(&isHost),
	})
}

func (pg *PostgresPlayerStore) FindByNickname(ctx context.Context, gameID int64, nickname string) (*Player, error) {
	player, err := pg.queries.FindPlayerByNickname(ctx, sqlc.FindPlayerByNicknameParams{
		GameID:   gameID,
		Nickname: nickname,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &Player{
		ID:       player.ID,
		Nickname: player.Nickname,
		IsHost:   fromPgBool(player.IsHost),
		GameID:   player.GameID,
		Status:   player.Status,
	}, nil
}

func (pg *PostgresPlayerStore) GetPlayerCountByGameCode(ctx context.Context, gameCode string) (int64, error) {
	return pg.queries.GetPlayerCountByGameCode(ctx, gameCode)
}
