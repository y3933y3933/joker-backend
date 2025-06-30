package api

import (
	"errors"
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/y3933y3933/joker/internal/service"
	"github.com/y3933y3933/joker/internal/store"
	"github.com/y3933y3933/joker/internal/utils/errx"
	"github.com/y3933y3933/joker/internal/utils/httpx"
	"github.com/y3933y3933/joker/internal/utils/param"
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
		msg, _ := ws.NewWSMessage(ws.MsgTypeGameStarted, ws.RoundStartedPayload{
			RoundID:          round.ID,
			QuestionPlayerID: round.QuestionPlayerID,
			AnswererID:       round.AnswerPlayerID,
		})
		room.Broadcast(msg)
	}

	httpx.SuccessResponse(c, round)
}

type SubmitQuestionRequest struct {
	QuestionID int64 `json:"questionID" binding:"required"`
}

func (h *RoundHandler) HandleSubmitQuestion(c *gin.Context) {
	var req SubmitQuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.BadRequestResponse(c, err)
		return
	}

	roundID, err := param.ParseIntParam(c, "id")
	if err != nil {
		httpx.BadRequestResponse(c, errors.New("invalid round id"))
		return
	}

	playerIDAny, ok := c.Get("player_id")
	if !ok {
		httpx.ServerErrorResponse(c, h.logger, errors.New("missing player id"))
		return
	}
	playerID := playerIDAny.(int64)

	err = h.roundService.SubmitQuestion(c.Request.Context(), roundID, req.QuestionID, playerID)
	if err != nil {
		switch {
		case errors.Is(err, errx.ErrForbidden):
			httpx.ForbiddenResponse(c, err)
		case errors.Is(err, errx.ErrInvalidStatus):
			httpx.BadRequestResponse(c, err)
		default:
			httpx.ServerErrorResponse(c, h.logger, err)
		}
		return
	}

	// 拿 round + question 資料（包含回答者 ID 與題目內容）
	round, err := h.roundService.GetRoundWithQuestion(c.Request.Context(), roundID)
	if err != nil {
		httpx.ServerErrorResponse(c, h.logger, err)
		return
	}

	gameAny, ok := c.Get("game")
	if !ok {
		httpx.ServerErrorResponse(c, h.logger, errors.New("missing game in context"))
		return
	}
	game := gameAny.(*store.Game)

	// 推播：給所有人通知已進入回答階段
	room := h.hub.GetRoom(game.Code)
	if room != nil {
		// 1️⃣ 推播給所有人：進入回答時間
		msg1, _ := ws.NewWSMessage(ws.MsgTypeAnswerTime, nil)
		room.Broadcast(msg1)

		// 2️⃣ 私訊給回答者：這是題目內容

		msg2, _ := ws.NewWSMessage(ws.MsgTypeRoundQuestion, map[string]string{
			"level":   round.Level,
			"content": round.Content,
		})

		room.SendTo(round.AnswerPlayerID, msg2)
	}

	httpx.SuccessResponse(c, gin.H{"message": "question submitted"})
}

type SubmitAnswerRequest struct {
	Answer string `json:"answer" binding:"required"`
}

func (h *RoundHandler) HandleSubmitAnswer(c *gin.Context) {
	var req SubmitAnswerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.BadRequestResponse(c, err)
		return
	}

	// 取得 roundID
	roundID, err := param.ParseIntParam(c, "id")
	if err != nil {
		httpx.BadRequestResponse(c, errors.New("invalid round id"))
		return
	}

	// 從 context 取得 playerID
	playerIDAny, ok := c.Get("player_id")
	if !ok {
		httpx.ServerErrorResponse(c, h.logger, errors.New("missing player id"))
		return
	}
	playerID := playerIDAny.(int64)

	// 呼叫 Service
	err = h.roundService.SubmitAnswer(c.Request.Context(), roundID, req.Answer, playerID)
	if err != nil {
		switch {
		case errors.Is(err, errx.ErrForbidden):
			httpx.ForbiddenResponse(c, err)
		case errors.Is(err, errx.ErrInvalidStatus):
			httpx.BadRequestResponse(c, err)
		case errors.Is(err, errx.ErrRoundNotFound):
			httpx.NotFoundResponse(c, err)
		default:
			httpx.ServerErrorResponse(c, h.logger, err)
		}
		return
	}

	// 取 game code 推播
	gameAny, ok := c.Get("game")
	if !ok {
		httpx.ServerErrorResponse(c, h.logger, errors.New("missing game context"))
		return
	}
	game := gameAny.(*store.Game)

	// 推播 answer_submitted 給所有人
	room := h.hub.GetRoom(game.Code)
	if room != nil {
		msg, _ := ws.NewWSMessage(ws.MsgTypeAnswerSubmitted, ws.AnswerSubmittedPayload{
			Answer: req.Answer,
		})
		room.Broadcast(msg)
	}

	httpx.SuccessResponse(c, gin.H{"message": "answer submitted"})
}

type DrawCardRequest struct {
	Index *int `json:"index" binding:"required"`
}

func (h *RoundHandler) HandleDrawCard(c *gin.Context) {
	var req DrawCardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.BadRequestResponse(c, err)
		return
	}

	roundID, err := param.ParseIntParam(c, "id")
	if err != nil {
		httpx.BadRequestResponse(c, errors.New("invalid round id"))
		return
	}

	playerIDAny, ok := c.Get("player_id")
	if !ok {
		httpx.ServerErrorResponse(c, h.logger, errors.New("missing player id"))
		return
	}
	playerID := playerIDAny.(int64)

	round, err := h.roundService.DrawCard(c.Request.Context(), roundID, playerID, *req.Index)
	if err != nil {
		switch {
		case errors.Is(err, errx.ErrForbidden):
			httpx.ForbiddenResponse(c, err)
		case errors.Is(err, errx.ErrInvalidStatus):
			httpx.BadRequestResponse(c, err)
		default:
			httpx.ServerErrorResponse(c, h.logger, err)
		}
		return
	}

	// 推播
	gameAny, _ := c.Get("game")
	game := gameAny.(*store.Game)

	room := h.hub.GetRoom(game.Code)
	if room != nil {
		if round.IsJoker {
			msg, _ := ws.NewWSMessage(ws.MsgTypeJokerRevealed, ws.JokerRevealedPayload{
				Level:   round.Level,
				Content: round.Content,
			})
			room.Broadcast(msg)
		} else {
			msg, _ := ws.NewWSMessage(ws.MsgTypePlayerSafe, nil)
			room.Broadcast(msg)
		}
	}

	httpx.SuccessResponse(c, gin.H{
		"joker": round.IsJoker,
	})
}

func (h *RoundHandler) HandleCreateNextRound(c *gin.Context) {
	gameAny, ok := c.Get("game")
	if !ok {
		httpx.ServerErrorResponse(c, h.logger, errors.New("missing game in context"))
		return
	}
	game := gameAny.(*store.Game)

	round, err := h.roundService.CreateNextRound(c.Request.Context(), game)
	if err != nil {
		httpx.ServerErrorResponse(c, h.logger, err)
		return
	}

	// 推播 round_started 給所有人
	room := h.hub.GetRoom(game.Code)
	if room != nil {
		msg, _ := ws.NewWSMessage(ws.MsgNextRoundStarted, ws.RoundStartedPayload{
			RoundID:          round.ID,
			AnswererID:       round.AnswerPlayerID,
			QuestionPlayerID: round.QuestionPlayerID,
		})
		room.Broadcast(msg)
	}

	httpx.SuccessResponse(c, round)
}
