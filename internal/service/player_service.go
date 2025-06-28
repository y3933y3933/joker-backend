package service

import (
	"context"

	"github.com/y3933y3933/joker/internal/store"
)

type PlayerService struct {
	playerStore store.PlayerStore
	gameStore   store.GameStore
}

func NewPlayerService(playerStore store.PlayerStore, gameStore store.GameStore) *PlayerService {
	return &PlayerService{
		playerStore: playerStore,
		gameStore:   gameStore,
	}
}

func (s *PlayerService) JoinGame(ctx context.Context, code, nickname string) (*store.Player, error) {
	game, err := s.gameStore.GetGameByCode(ctx, code)
	if err != nil {
		return nil, err
	}
	count, err := s.playerStore.CountPlayerInGame(ctx, game.ID)
	if err != nil {
		return nil, err
	}

	isHost := count == 0
	args := &store.Player{
		Nickname: nickname,
		IsHost:   isHost,
		GameID:   game.ID,
	}
	player, err := s.playerStore.Create(ctx, args)

	if err != nil {
		return nil, err
	}

	return player, nil

}
