package api

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/y3933y3933/joker/internal/service"
)

type GameHandler struct {
	gameService *service.GameService
	logger      *slog.Logger
}

func NewGameHandler(gameService *service.GameService, logger *slog.Logger) *GameHandler {
	return &GameHandler{
		gameService: gameService,
		logger:      logger,
	}
}

func (h *GameHandler) HandleCreateGame(c *gin.Context) {
	game, err := h.gameService.CreateGame(c.Request.Context())
	if err != nil {
		ServerErrorResponse(c, h.logger, err)
		return
	}
	SuccessResponse(c, game)
}
