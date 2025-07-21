package service

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/y3933y3933/joker/internal/store"
	"github.com/y3933y3933/joker/internal/utils/codegen"
	"github.com/y3933y3933/joker/internal/utils/errx"
)

type GameService struct {
	gameStore   store.GameStore
	playerStore store.PlayerStore
}

func NewGameService(gameStore store.GameStore, playerStore store.PlayerStore) *GameService {
	return &GameService{
		gameStore:   gameStore,
		playerStore: playerStore,
	}
}

func (s *GameService) generateCode(ctx context.Context) (string, error) {
	const maxTries = 10

	for i := 0; i < maxTries; i++ {
		code := codegen.GenerateCode(6)
		exists, err := s.gameStore.GameCodeExists(ctx, code)
		if err != nil {
			return "", fmt.Errorf("check game code: %w", err)
		}
		if !exists {
			return code, nil
		}
	}

	return "", errx.ErrGenerateCode
}

func (s *GameService) CreateGame(ctx context.Context) (*store.Game, error) {
	code, err := s.generateCode(ctx)
	if err != nil {
		return nil, err
	}

	args := &store.Game{
		Code:   code,
		Status: store.GameStatusWaiting,
	}
	game, err := s.gameStore.Create(ctx, args)
	if err != nil {
		return nil, err
	}

	return game, nil
}

func (s *GameService) EndGame(ctx context.Context, code string) error {
	game, err := s.gameStore.GetGameByCode(ctx, code)
	if err != nil {
		return err
	}

	if game.Status == store.GameStatusEnded {
		return errx.ErrInvalidGameStatus
	}

	return s.gameStore.EndGame(ctx, code)
}

func (s *GameService) GetGameSummaryByCode(ctx context.Context, gameID int64) (*store.GameSummary, error) {
	stats, err := s.gameStore.GetGameSummary(ctx, gameID)
	if err != nil {
		return nil, err
	}

	playerStats, err := s.gameStore.GetGamePlayerStats(ctx, gameID)
	if err != nil {
		return nil, err
	}

	return &store.GameSummary{
		TotalRounds: stats.TotalRounds,
		JokerCards:  stats.JokerCards,
		Players:     playerStats,
	}, nil
}

func (s *GameService) DeleteGameIfEmpty(ctx context.Context, gameCode string) error {
	count, err := s.playerStore.GetPlayerCountByGameCode(ctx, gameCode)
	if err != nil {
		return err
	}

	if count > 0 {
		return nil
	}

	return s.gameStore.DeleteByCode(ctx, gameCode)
}

func (s *GameService) GetGameByCode(ctx context.Context, gameCode string) (*store.Game, error) {
	return s.gameStore.GetGameByCode(ctx, gameCode)
}

type GameQueryParams struct {
	Code     string `json:"code"`
	Status   string `json:"status"`
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
}

func (s *GameService) ListGame(ctx context.Context, query GameQueryParams) (*store.PaginatedGame, error) {
	filters := store.Filters{
		Page:     query.Page,
		PageSize: query.PageSize,
	}

	result, err := s.gameStore.List(ctx, query.Code, query.Status, filters)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *GameService) ValidateGameParams(params GameQueryParams) error {
	if params.Status != "" {
		if !slices.Contains([]string{"waiting", "playing", "ended"}, params.Status) {
			return errors.New("invalid status: must be 'waiting', 'playing', or 'ended'")
		}

	}

	if params.Page < 1 {
		return errors.New("page must be greater than 0")
	}

	if params.PageSize < 1 || params.PageSize > 100 {
		return errors.New("page_size must be between 1 and 100")
	}

	return nil
}
