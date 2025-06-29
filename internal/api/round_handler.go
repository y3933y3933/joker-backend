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

type RoundHandler struct {
	roundService *service.RoundService
	hub          *ws.Hub
	logger       *slog.Logger
}

func NewRoundHandler(roundService *service.RoundService, logger *slog.Logger, hub *ws.Hub) *RoundHandler {
	return &RoundHandler{
		roundService: roundService,
		hub:          hub,
		logger:       logger,
	}
}

func (h *RoundHandler) HandleStartGame(c *gin.Context) {
	gameAny, ok := c.Get("game")
	if !ok {
		httpx.ServerErrorResponse(c, h.logger, errors.New("missing game in context"))
		return
	}
	game := gameAny.(*store.Game)

	round, err := h.roundService.StartGame(c.Request.Context(), game)
	if err != nil {
		switch {
		case errors.Is(err, errx.ErrInvalidGameStatus):
			httpx.BadRequestResponse(c, errors.New("game already started or ended"))

		case errors.Is(err, errx.ErrInvalidGameStatus):
			httpx.BadRequestResponse(c, errors.New("game already started or ended"))

		default:
			httpx.ServerErrorResponse(c, h.logger, err)
		}
		return

	}

	// ✅ 推播給所有人
	room := h.hub.GetRoom(game.Code)
	if room != nil {
		msg, _ := ws.NewWSMessage(ws.MsgTypeGameStarted, ws.GameStartedPayload{
			RoundID:          round.ID,
			QuestionPlayerID: round.QuestionPlayerID,
			AnswererID:       round.AnswerPlayerID,
		})
		room.Broadcast(msg)
	}

	httpx.SuccessResponse(c, round)
}
