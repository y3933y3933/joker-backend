package api

import (
	"log/slog"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/y3933y3933/joker/internal/service"
	"github.com/y3933y3933/joker/internal/utils/httpx"
)

type GameHandler struct {
	gameService     *service.GameService
	questionService *service.QuestionService
	logger          *slog.Logger
}

func NewGameHandler(gameService *service.GameService, questionService *service.QuestionService, logger *slog.Logger) *GameHandler {
	return &GameHandler{
		gameService:     gameService,
		questionService: questionService,
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
	limitStr := c.DefaultQuery("limit", "5")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 5
	}

	questions, err := h.questionService.ListRandomQuestions(c.Request.Context(), limit)
	if err != nil {
		httpx.ServerErrorResponse(c, h.logger, err)
		return
	}

	httpx.SuccessResponse(c, questions)
}
