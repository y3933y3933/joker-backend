package api

import (
	"errors"
	"log/slog"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/y3933y3933/joker/internal/service"
	"github.com/y3933y3933/joker/internal/store"
	"github.com/y3933y3933/joker/internal/utils/errx"
	"github.com/y3933y3933/joker/internal/utils/httpx"
	"github.com/y3933y3933/joker/internal/ws"
)

type GameHandler struct {
	gameService     *service.GameService
	questionService *service.QuestionService
	hub             *ws.Hub
	logger          *slog.Logger
}

func NewGameHandler(gameService *service.GameService, questionService *service.QuestionService, hub *ws.Hub, logger *slog.Logger) *GameHandler {
	return &GameHandler{
		gameService:     gameService,
		questionService: questionService,
		hub:             hub,
		logger:          logger,
	}
}

func (h *GameHandler) HandleCreateGame(c *gin.Context) {
	game, err := h.gameService.CreateGame(c.Request.Context())
	if err != nil {
		httpx.ServerErrorResponse(c, h.logger, err)
		return
	}
	httpx.SuccessResponse(c, game)
}

func (h *GameHandler) HandleGetQuestions(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "3")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 3
	}

	questions, err := h.questionService.ListRandomQuestions(c.Request.Context(), limit)
	if err != nil {
		httpx.ServerErrorResponse(c, h.logger, err)
		return
	}

	httpx.SuccessResponse(c, questions)
}

func (h *GameHandler) HandleEndGame(c *gin.Context) {
	gameAny, exists := c.Get("game")
	if !exists {
		httpx.ServerErrorResponse(c, h.logger, errors.New("missing game in context"))
		return
	}
	game := gameAny.(*store.Game)

	err := h.gameService.EndGame(c.Request.Context(), game.Code, game.Status)
	if err != nil {
		if errors.Is(err, errx.ErrInvalidGameStatus) {
			httpx.BadRequestResponse(c, err)
			return
		}
		httpx.ServerErrorResponse(c, h.logger, err)
		return
	}

	// 推播 game_ended 給所有人（若有 hub）
	if room := h.hub.GetRoom(game.Code); room != nil {
		msg, _ := ws.NewWSMessage(ws.MsgTypeGameEnded, gin.H{"gameCode": game.Code})
		room.Broadcast(msg)
	}

	httpx.SuccessResponse(c, gin.H{"message": "game ended"})
}

func (h *GameHandler) GetGameSummary(c *gin.Context) {
	gameAny, exists := c.Get("game")
	if !exists {
		httpx.ServerErrorResponse(c, h.logger, errors.New("missing game in context"))
		return
	}
	game := gameAny.(*store.Game)

	summary, err := h.gameService.GetGameSummaryByCode(c.Request.Context(), game.ID)
	if err != nil {
		httpx.ServerErrorResponse(c, h.logger, errors.New("failed to get game summary"))
		return
	}
	httpx.SuccessResponse(c, summary)
}
