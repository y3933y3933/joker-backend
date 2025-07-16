package service

import (
	"context"

	"github.com/y3933y3933/joker/internal/store"
)

type UserService struct {
	userStore store.UserStore
}

func NewUserService(userStore store.UserStore) *UserService {
	return &UserService{userStore: userStore}
}

func (s *UserService) GetUserInfo(ctx context.Context, userID int64) (*store.User, error) {
	return s.userStore.GetUserByID(ctx, userID)
}
