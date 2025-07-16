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
	"github.com/y3933y3933/joker/internal/utils/errx"
	"github.com/y3933y3933/joker/internal/utils/httpx"
)

type Middleware struct {
	gameService *service.GameService
	authService *service.AuthService
}

func NewMiddleware(gameService *service.GameService,
	authService *service.AuthService) *Middleware {
	return &Middleware{
		gameService: gameService,
		authService: authService,
	}
}

type contextKey string

var codePattern = regexp.MustCompile(`^[A-Za-z0-9]{6}$`)

func (m *Middleware) ValidateGameExists() gin.HandlerFunc {
	return func(c *gin.Context) {
		code := c.Param("code")
		if code == "" || !codePattern.MatchString(code) {
			httpx.BadRequestResponse(c, errors.New("invalid game code"))
			return
		}

		game, err := m.gameService.GetGameByCode(c.Request.Context(), code)

		if err != nil {
			httpx.NotFoundResponse(c, errors.New("game not found"))
			return
		}

		c.Set("game", game)
		c.Next()
	}
}

func (m *Middleware) WithPlayerID() gin.HandlerFunc {
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

func (m *Middleware) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Header("Vary", "Authorization")

		tokenString := extractTokenFromHeaders(c.Request.Header)
		if tokenString == "" {
			httpx.UnAuthorized(c, errx.ErrInvalidAuthorizationHeader)
			return
		}

		claims, err := m.authService.ParseToken(tokenString)
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

		c.Set("user_id", claims.UserID)

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

func (m *Middleware) RequireUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, exists := c.Get("user_id")
		if !exists {
			httpx.UnAuthorized(c, errx.ErrLoginRequired)
			return
		}
		c.Next()
	}
}
