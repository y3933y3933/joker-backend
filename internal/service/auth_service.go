package service

import (
	"context"

	"github.com/y3933y3933/joker/internal/store"
)

type AuthService struct {
	userStore store.UserStore
}

func NewAuthService(userStore store.UserStore) *AuthService {
	return &AuthService{userStore: userStore}
}

func (s *AuthService) CreateUser(ctx context.Context, username, password string) (*store.User, error) {
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
