package service

import (
	"context"

	"github.com/y3933y3933/joker/internal/store"
)

type UserService struct {
	userStore store.UserStore
}

func NewUserStore(userStore store.UserStore) *UserService {
	return &UserService{userStore: userStore}
}

func (s *UserService) CreateUser(ctx context.Context, username, password string) (*store.User, error) {
	user := &store.User{
		Username: username,
	}
	err := user.Password.Set(password)
	if err != nil {
		return nil, err
	}

	err = s.userStore.Create(ctx, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}
