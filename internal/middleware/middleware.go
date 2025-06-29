package middleware

import (
	"errors"
	"regexp"
	"strconv"

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

func WithPlayerID() gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Query("player_id") // 可改成從 token 或 header 取出
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			httpx.BadRequestResponse(c, errors.New("invalid player id"))
			return
		}
		c.Set("player_id", id)
		c.Next()
	}
}
