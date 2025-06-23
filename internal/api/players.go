package api

import (
	"database/sql"
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/y3933y3933/joker/internal/database"
	"github.com/y3933y3933/joker/internal/utils"
	"github.com/y3933y3933/joker/internal/ws"
)

type PlayersHandler struct {
	logger  *slog.Logger
	queries *database.Queries
	hub     *ws.Hub
}

func NewPlayersHandler(queries *database.Queries, logger *slog.Logger, hub *ws.Hub) *PlayersHandler {
	return &PlayersHandler{
		logger:  logger,
		queries: queries,
		hub:     hub,
	}
}

type PlayerResponse struct {
	ID       int64  `json:"id"`
	Nickname string `json:"nickname"`
	IsHost   bool   `json:"isHost"`
}

func (h *PlayersHandler) ListPlayers(c *gin.Context) {
	ctx := c.Request.Context()
	code := c.Param("code")

	players, err := h.queries.ListPlayersByGameCode(ctx, code)
	if err != nil {
		h.logger.Error("list players by game code error: ", err)
		switch {
		case errors.Is(err, sql.ErrNoRows):
			BadRequest(c, "game not found")
		default:
			InternalServerError(c, "something went wrong")
		}
		return
	}

	Success(c, transformToPlayerResponse(players))
}

func transformToPlayerResponse(players []database.ListPlayersByGameCodeRow) []PlayerResponse {
	var playerResponses []PlayerResponse
	for _, p := range players {
		playerResponses = append(playerResponses, PlayerResponse{
			ID:       p.ID,
			Nickname: p.Nickname,
			IsHost:   p.IsHost.Bool,
		})
	}
	return playerResponses
}

type JoinGameRequest struct {
	Nickname string `json:"nickname" binding:"required"`
}

func (h *PlayersHandler) JoinGame(c *gin.Context) {
	ctx := c.Request.Context()
	gameCode := c.Param("code")

	var req JoinGameRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "Invalid nickname")
		return
	}

	game, err := h.queries.GetGameByCode(ctx, gameCode)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			NotFound(c, "game not found")
		} else {
			InternalServerError(c, "DB error")
		}
		return
	}

	count, err := h.queries.CountPlayersInGame(ctx, game.ID)
	if err != nil {
		InternalServerError(c, "Count error")
		return
	}
	isHost := count == 0

	player, err := h.queries.CreatePlayer(ctx, database.CreatePlayerParams{
		GameID:   game.ID,
		Nickname: req.Nickname,
		IsHost:   pgtype.Bool{Bool: isHost, Valid: true},
	})
	if err != nil {
		InternalServerError(c, "Create player failed")
		return
	}

	// ✅ WebSocket 廣播
	h.hub.BroadcastToGame(game.Code, ws.WebSocketMessage{
		Type: "player_joined",
		Data: gin.H{
			"id":       player.ID,
			"nickname": player.Nickname,
			"isHost":   player.IsHost,
		},
	})

	// ✅ 回傳該玩家資訊
	Success(c, PlayerResponse{
		ID:       player.ID,
		Nickname: player.Nickname,
		IsHost:   player.IsHost.Bool,
	})

}

func (h *RoundsHandler) RemovePlayer(c *gin.Context) {
	ctx := c.Request.Context()
	gameCode := c.Param("code")
	playerIDParam := c.Param("player_id")

	playerID, err := utils.ParseID(playerIDParam)
	if err != nil {
		BadRequest(c, "invalid player ID")
		return
	}

	game, err := h.queries.GetGameByCode(ctx, gameCode)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			NotFound(c, "game not found")
			return
		}
		InternalServerError(c, "db error")
		return
	}

	player, err := h.queries.DeletePlayer(ctx, database.DeletePlayerParams{
		ID:     playerID,
		GameID: game.ID,
	})
	if err != nil {
		InternalServerError(c, "failed to remove player")
		return
	}

	// 廣播玩家離開
	h.hub.BroadcastToGame(game.Code, ws.WebSocketMessage{
		Type: "player_left",
		Data: gin.H{
			"id":       player.ID,
			"nickname": player.Nickname,
		},
	})

	c.Status(http.StatusNoContent)
}
