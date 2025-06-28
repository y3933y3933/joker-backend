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

func (s *PlayerService) JoinGame(ctx context.Context, gameID int64, nickname string) (*store.Player, error) {
	count, err := s.playerStore.CountPlayerInGame(ctx, gameID)
	if err != nil {
		return nil, err
	}

	isHost := count == 0
	args := &store.Player{
		Nickname: nickname,
		IsHost:   isHost,
		GameID:   gameID,
	}
	player, err := s.playerStore.Create(ctx, args)

	if err != nil {
		return nil, err
	}

	return player, nil

}

func (s *PlayerService) ListPlayersInGame(ctx context.Context, gameID int64) ([]*store.Player, error) {
	return s.playerStore.FindPlayersByGameID(ctx, gameID)
}
