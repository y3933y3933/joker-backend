package service

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/y3933y3933/joker/internal/store"
	"github.com/y3933y3933/joker/internal/utils/errx"
)

type CustomClaims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	// Role     string `json:"role"`
	jwt.RegisteredClaims
}

type AuthService struct {
	userStore store.UserStore
	jwtSecret []byte
}

func NewAuthService(userStore store.UserStore, jwtSecret []byte) *AuthService {
	return &AuthService{userStore: userStore, jwtSecret: jwtSecret}
}

func (s *AuthService) CreateUser(ctx context.Context, username, password string) (*store.User, error) {
	user := &store.User{
		Username: username,
	}
	err := user.Password.Set(password)
	if err != nil {
		return nil, err
	}

	id, err := s.userStore.Create(ctx, user)
	if err != nil {
		return nil, err
	}
	user.ID = id
	return user, nil
}

func (s *AuthService) Login(ctx context.Context, username, password string) (string, error) {
	user, err := s.userStore.GetUserByUsername(ctx, username)
	if err != nil {
		return "", err
	}
	passwordIsMatch, err := user.Password.Matches(password)
	if err != nil {
		return "", err
	}

	if !passwordIsMatch {
		return "", errx.ErrInvalidCredentials
	}

	token, err := s.createToken(user.ID, user.Username)
	if err != nil {
		return "", err
	}

	return token, nil // check & createToken
}

func (s *AuthService) createToken(userID int64, username string) (string, error) {
	now := time.Now()

	claims := CustomClaims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(now),
			// Subject:  userID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(s.jwtSecret)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *AuthService) ParseToken(tokenString string) (*CustomClaims, error) {

	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return s.jwtSecret, nil
	})

	if err != nil {
		return nil, errx.ErrInvalidToken
	}

	if !token.Valid {
		return nil, errx.ErrInvalidToken
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return nil, errx.ErrInvalidToken
	}

	return claims, nil

}
