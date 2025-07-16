package service

import (
	"context"

	"github.com/y3933y3933/joker/internal/store"
)

type AdminService struct {
	playerStore   store.PlayerStore
	feedbackStore store.FeedbackStore
	gameStore     store.GameStore
}

func NewAdminService(playerStore store.PlayerStore, feedbackStore store.FeedbackStore, gameStore store.GameStore) *AdminService {
	return &AdminService{
		playerStore:   playerStore,
		feedbackStore: feedbackStore,
		gameStore:     gameStore,
	}
}

type DashboardData struct {
	GamesTodayCount       int64 `json:"gamesTodayCount"`
	ActiveRoomsCount      int64 `json:"activeRoomsCount"`
	FeedbackOneMonthCount int64 `json:"feedbackOneMonthCount"`
	LivePlayerCount       int64 `json:"livePlayerCount"`
}

func (s *AdminService) GetDashboardData(ctx context.Context) (*DashboardData, error) {
	activeRoomsCount, err := s.gameStore.GetActiveRoomsCount(ctx)
	if err != nil {
		return nil, err
	}

	gamesTodayCount, err := s.gameStore.GetGamesTodayCount(ctx)
	if err != nil {
		return nil, err
	}
	feedbackOneMonthCount, err := s.feedbackStore.CountRecentFeedbacksOneMonth(ctx)
	if err != nil {
		return nil, err
	}
	livePlayerCount, err := s.playerStore.GetLivePlayerCount(ctx)
	if err != nil {
		return nil, err
	}

	return &DashboardData{
		ActiveRoomsCount:      activeRoomsCount,
		GamesTodayCount:       gamesTodayCount,
		FeedbackOneMonthCount: feedbackOneMonthCount,
		LivePlayerCount:       livePlayerCount,
	}, nil

}
