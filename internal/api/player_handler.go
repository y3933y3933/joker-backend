package api

import (
	"errors"
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/y3933y3933/joker/internal/service"
	"github.com/y3933y3933/joker/internal/store"
	"github.com/y3933y3933/joker/internal/utils/errx"
)

type PlayerHandler struct {
	playerService service.PlayerService
	logger        *slog.Logger
}

func NewPlayerHandler(playerService service.PlayerService, logger *slog.Logger) *PlayerHandler {
	return &PlayerHandler{
		playerService: playerService,
		logger:        logger,
	}
}

type JoinGameRequest struct {
	Nickname string `json:"nickname" binding:"required"`
}

func (h *PlayerHandler) HandleJoinGame(c *gin.Context) {
	var req JoinGameRequest
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		BadRequestResponse(c, err)
		return
	}

	gameAny, ok := c.Get("game")
	if !ok {
		ServerErrorResponse(c, h.logger, errors.New("missing game in context"))
		return
	}

	game := gameAny.(*store.Game)

	player, err := h.playerService.JoinGame(c.Request.Context(), game.Code, req.Nickname)
	if err != nil {
		if errors.Is(err, errx.ErrGameNotFound) {
			BadRequestResponse(c, errors.New("game not found"))
			return
		}
		ServerErrorResponse(c, h.logger, err)
		return
	}

	SuccessResponse(c, player)

}
