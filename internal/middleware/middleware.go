package middleware

import (
	"errors"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/y3933y3933/joker/internal/store"
	"github.com/y3933y3933/joker/internal/utils/httpx"
)

var codePattern = regexp.MustCompile(`^[A-Za-z0-9]{6}$`)

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
