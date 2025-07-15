package middleware

import (
	"errors"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/y3933y3933/joker/internal/service"
	"github.com/y3933y3933/joker/internal/store"
	"github.com/y3933y3933/joker/internal/utils/errx"
	"github.com/y3933y3933/joker/internal/utils/httpx"
)

var codePattern = regexp.MustCompile(`^[A-Za-z0-9]{6}$`)

// TODO: replace gameStore
func ValidateGameExists(gameStore store.GameStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		code := c.Param("code")
		if code == "" || !codePattern.MatchString(code) {
			httpx.BadRequestResponse(c, errors.New("invalid game code"))
			return
		}

		game, err := gameStore.GetGameByCode(c.Request.Context(), code)
		if err != nil {
			httpx.NotFoundResponse(c, errors.New("game not found"))
			return
		}

		c.Set("game", game) // 讓後續 handler 可用
		c.Next()
	}
}

func WithPlayerID() gin.HandlerFunc {
	return func(c *gin.Context) {
		playerIDStr := c.GetHeader("X-Player-ID")
		if playerIDStr == "" {
			httpx.BadRequestResponse(c, errors.New("missing X-Player-ID header"))
			return
		}
		playerID, err := strconv.ParseInt(playerIDStr, 10, 64)
		if err != nil {
			httpx.BadRequestResponse(c, errors.New("invalid player id"))
			return
		}
		c.Set("player_id", playerID)
		c.Next()
	}
}

type AuthMiddleware struct {
	AuthService service.AuthService
}

type contextKey string

const UserContextKey = contextKey("user")

func (m *AuthMiddleware) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Header("Vary", "Authorization")

		tokenString := extractTokenFromHeaders(c.Request.Header)
		if tokenString == "" {
			httpx.UnAuthorized(c, errx.ErrInvalidAuthorizationHeader)
			return
		}

		claims, err := m.AuthService.ParseToken(tokenString)
		if err != nil {
			httpx.UnAuthorized(c, err)
			return
		}

		if claims.ExpiresAt == nil {
			httpx.UnAuthorized(c, errx.ErrInvalidToken)
			return
		}

		if time.Now().After(claims.ExpiresAt.Time) {
			httpx.UnAuthorized(c, errx.ErrTokenExpired)
			return
		}

		c.Set("userID", claims.UserID)

		c.Next()
	}
}

func extractTokenFromHeaders(headers http.Header) string {
	authHeader := headers.Get("Authorization")

	if authHeader == "" {
		return ""
	}

	headerParts := strings.Split(authHeader, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return ""
	}
	return headerParts[1]
}
