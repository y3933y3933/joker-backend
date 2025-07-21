package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

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

type GamePlayerSummary struct {
	ID              int64  `json:"id"`
	Nickname        string `json:"nickname"`
	JokerCardsDrawn int32  `json:"jokerCardsDrawn"`
}

type GameSummary struct {
	TotalRounds int64               `json:"totalRounds"`
	JokerCards  int64               `json:"jokerCards"`
	Players     []GamePlayerSummary `json:"players"`
}

type AdminGame struct {
	ID          int64     `json:"id"`
	Code        string    `json:"code"`
	Status      string    `json:"status"`
	PlayerCount int64     `json:"playerCount"`
	CreatedAt   time.Time `json:"createdAt"`
}

type PaginatedGame struct {
	Games []AdminGame `json:"games"`
	Metadata
}

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
	GetGameSummary(ctx context.Context, gameID int64) (*GameSummary, error)
	GetGamePlayerStats(ctx context.Context, gameID int64) ([]GamePlayerSummary, error)
	GetGameStatusByID(ctx context.Context, gameID int64) (string, error)
	DeleteByCode(ctx context.Context, gameCode string) error
	GetGamesTodayCount(ctx context.Context) (int64, error)
	GetActiveRoomsCount(ctx context.Context) (int64, error)
	List(ctx context.Context, code, status string, filters Filters) (*PaginatedGame, error)
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

func (pg *PostgresGameStore) GetGameSummary(ctx context.Context, gameID int64) (*GameSummary, error) {
	summary, err := pg.queries.GetGameSummaryStats(ctx, gameID)

	return &GameSummary{
		TotalRounds: summary.TotalRounds,
		JokerCards:  summary.JokerCards,
	}, err
}

func (pg *PostgresGameStore) GetGamePlayerStats(ctx context.Context, gameID int64) ([]GamePlayerSummary, error) {
	dbPlayers, err := pg.queries.GetGamePlayerStats(ctx, gameID)
	if err != nil {
		return nil, err
	}

	players := make([]GamePlayerSummary, len(dbPlayers))

	for i, p := range dbPlayers {
		players[i] = GamePlayerSummary{
			ID:              p.ID,
			Nickname:        p.Nickname,
			JokerCardsDrawn: int32(p.JokerCardsDrawn),
		}
	}
	return players, nil
}

func (pg *PostgresGameStore) GetGameStatusByID(ctx context.Context, gameID int64) (string, error) {
	status, err := pg.queries.GetGameStatusByID(ctx, gameID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", errx.ErrGameNotFound
		}
		return "", err
	}

	return status, nil
}

func (pg *PostgresGameStore) DeleteByCode(ctx context.Context, gameCode string) error {
	return pg.queries.DeleteByCode(ctx, gameCode)
}

func (pg *PostgresGameStore) GetGamesTodayCount(ctx context.Context) (int64, error) {
	return pg.queries.GetGamesTodayCount(ctx)
}

func (pg *PostgresGameStore) GetActiveRoomsCount(ctx context.Context) (int64, error) {
	return pg.queries.GetActiveRoomsCount(ctx)
}

func (pg *PostgresGameStore) List(ctx context.Context, code, status string, filters Filters) (*PaginatedGame, error) {

	args := sqlc.ListGamesParams{
		Upper:  code,
		Status: status,
		Limit:  int32(filters.limit()),
		Offset: int32(filters.offset()),
	}

	fmt.Printf("Listing games with params: %+v\n", args)
	rows, err := pg.queries.ListGames(ctx, args)

	if err != nil {
		return nil, err
	}

	var totalCount = 0
	gameResponse := make([]AdminGame, len(rows))
	for i, f := range rows {
		gameResponse[i] = AdminGame{
			ID:          f.ID,
			Code:        f.Code,
			Status:      f.Status,
			CreatedAt:   f.CreatedAt.Time,
			PlayerCount: f.PlayerCount,
		}
		totalCount = int(f.TotalCount)
	}

	return &PaginatedGame{
		Games:    gameResponse,
		Metadata: CalculateMetadata(totalCount, filters.Page, filters.PageSize),
	}, nil
}
