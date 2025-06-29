package service

import (
	"context"
	"fmt"

	"github.com/y3933y3933/joker/internal/store"
	"github.com/y3933y3933/joker/internal/utils/codegen"
	"github.com/y3933y3933/joker/internal/utils/errx"
)

type GameService struct {
	gameStore store.GameStore
}

func NewGameService(gameStore store.GameStore) *GameService {
	return &GameService{
		gameStore: gameStore,
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

func (s *GameService) EndGame(ctx context.Context, code string, status string) error {
	if status == store.GameStatusEnded {
		return errx.ErrInvalidGameStatus
	}

	return s.gameStore.EndGame(ctx, code)
}
