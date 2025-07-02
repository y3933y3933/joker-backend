package service

import (
	"context"

	"github.com/y3933y3933/joker/internal/store"
	"github.com/y3933y3933/joker/internal/utils/errx"
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
	// 🔍 檢查暱稱是否已存在
	existing, err := s.playerStore.FindByNickname(ctx, gameID, nickname)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errx.ErrDuplicateNickname
	}

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

func (s *PlayerService) LeaveGame(ctx context.Context, playerID int64) (left *store.Player, newHost *store.Player, err error) {
	player, err := s.playerStore.FindByID(ctx, playerID)
	if err != nil {
		return nil, nil, err
	}

	gameStatus, err := s.gameStore.GetGameStatusByID(ctx, player.GameID)
	if err != nil {
		return nil, nil, err
	}

	if gameStatus != store.GameStatusWaiting {
		return nil, nil, errx.ErrGameAlreadyStarted
	}

	err = s.playerStore.DeleteByID(ctx, playerID)
	if err != nil {
		return nil, nil, err
	}

	if player.IsHost {
		newHost, err := s.TransferHost(ctx, player)
		if err != nil {
			return nil, nil, err
		}
		return player, newHost, nil

	}

	return player, nil, nil
}

func (s *PlayerService) TransferHost(ctx context.Context, player *store.Player) (*store.Player, error) {
	players, err := s.playerStore.FindOnlinePlayersByGameID(ctx, player.GameID)
	if err != nil {
		return nil, err
	}
	if len(players) > 0 {
		newHost := players[0]
		err = s.playerStore.UpdateHost(ctx, newHost.ID, true)
		if err != nil {
			return nil, err
		}
		return newHost, nil
	}
	return nil, errx.ErrNotEnoughPlayers
}

func (s *PlayerService) MarkPlayerDisconnected(ctx context.Context, playerID int64) error {
	return s.playerStore.UpdatePlayerStatus(ctx, playerID, store.PlayerStatusOffline)
}

func (s *PlayerService) FindPlayerByID(ctx context.Context, playerID int64) (*store.Player, error) {
	return s.playerStore.FindByID(ctx, playerID)
}
