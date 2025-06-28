package api

import (
	"errors"
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/y3933y3933/joker/internal/service"
	"github.com/y3933y3933/joker/internal/store"
	"github.com/y3933y3933/joker/internal/utils/errx"
	"github.com/y3933y3933/joker/internal/utils/httpx"
	"github.com/y3933y3933/joker/internal/ws"
)

type PlayerHandler struct {
	playerService service.PlayerService
	hub           *ws.Hub
	logger        *slog.Logger
}

func NewPlayerHandler(playerService service.PlayerService, hub *ws.Hub, logger *slog.Logger) *PlayerHandler {
	return &PlayerHandler{
		playerService: playerService,
		hub:           hub,
		logger:        logger,
	}
}

type JoinGameRequest struct {
	Nickname string `json:"nickname" binding:"required"`
}

func (h *PlayerHandler) HandleJoinGame(c *gin.Context) {
	var req JoinGameRequest
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		httpx.BadRequestResponse(c, err)
		return
	}

	gameAny, ok := c.Get("game")
	if !ok {
		httpx.ServerErrorResponse(c, h.logger, errors.New("missing game in context"))
		return
	}

	game := gameAny.(*store.Game)

	player, err := h.playerService.JoinGame(c.Request.Context(), game.ID, req.Nickname)
	if err != nil {
		if errors.Is(err, errx.ErrGameNotFound) {
			httpx.BadRequestResponse(c, errors.New("game not found"))
			return
		}
		httpx.ServerErrorResponse(c, h.logger, err)
		return
	}

	// ✅ 推播 player_joined 給房間內所有人
	room := h.hub.GetRoom(game.Code)
	if room != nil {
		msg, err := ws.NewWSMessage("player_joined", ws.PlayerJoinedPayload{
			ID:       player.ID,
			Nickname: player.Nickname,
			IsHost:   player.IsHost,
		})
		if err != nil {
			httpx.ServerErrorResponse(c, h.logger, err)
			return
		}
		room.Broadcast(msg)
	}

	httpx.SuccessResponse(c, player)

}

func (h *PlayerHandler) HandleListPlayers(c *gin.Context) {
	gameAny, ok := c.Get("game")
	if !ok {
		httpx.ServerErrorResponse(c, h.logger, errors.New("missing game in context"))
		return
	}
	game := gameAny.(*store.Game)

	players, err := h.playerService.ListPlayersInGame(c.Request.Context(), game.ID)
	if err != nil {
		httpx.ServerErrorResponse(c, h.logger, err)
		return
	}

	httpx.SuccessResponse(c, players)
}
