package api

import (
	"database/sql"
	"errors"
	"log/slog"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/y3933y3933/joker/internal/database"
	"github.com/y3933y3933/joker/internal/utils"
	"github.com/y3933y3933/joker/internal/ws"
)

type RoundsHandler struct {
	logger  *slog.Logger
	queries *database.Queries
	hub     *ws.Hub
}

func NewRoundsHandler(queries *database.Queries, logger *slog.Logger, hub *ws.Hub) *RoundsHandler {
	return &RoundsHandler{
		logger:  logger,
		queries: queries,
		hub:     hub,
	}
}

func (h *RoundsHandler) StartGame(c *gin.Context) {
	ctx := c.Request.Context()
	code := c.Param("code")

	game, err := h.queries.GetGameByCode(ctx, code)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			NotFound(c, "game not found")
			return
		}
		InternalServerError(c, "get game failed")
		return
	}

	// 取得所有玩家（照加入順序）
	players, err := h.queries.ListPlayersByGameCode(ctx, game.Code)
	if err != nil || len(players) < 2 {
		InternalServerError(c, "need at least 2 players")
		return
	}

	// 找出出題者與回答者
	questioner := players[0]
	answerer := players[1%len(players)] // 支援兩人時也 OK

	deck := utils.CreateShuffledDeck(1, 3)

	// 建立 round（暫不指定題目）
	round, err := h.queries.CreateRound(ctx, database.CreateRoundParams{
		GameID:           game.ID,
		QuestionID:       pgtype.Int8{Valid: false},
		QuestionPlayerID: questioner.ID,
		AnswerPlayerID:   pgtype.Int8{Int64: answerer.ID, Valid: true},
		Deck:             deck,
	})
	if err != nil {
		InternalServerError(c, "create round failed")
		return
	}

	// 廣播開始遊戲 & 告知出題者與回答者
	h.hub.BroadcastToGame(game.Code, ws.WebSocketMessage{
		Type: "game_started",
		Data: gin.H{
			"roundId":      round.ID,
			"questionerId": questioner.ID,
			"answererId":   answerer.ID,
			"status":       round.Status,
		},
	})

	Success(c, gin.H{
		"roundId":      round.ID,
		"questionerId": questioner.ID,
		"answererId":   answerer.ID,
	})
}

type SubmitQuestionRequest struct {
	QuestionID int64 `json:"questionId" binding:"required"`
}

func (h *RoundsHandler) SubmitQuestion(c *gin.Context) {
	ctx := c.Request.Context()
	code := c.Param("code")
	roundIDStr := c.Param("id")
	roundID, err := utils.ParseID(roundIDStr)
	if err != nil {
		BadRequest(c, "invalid round id")
		return
	}

	var req SubmitQuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "invalid questionId")
		return
	}

	game, err := h.queries.GetGameByCode(ctx, code)
	if err != nil {
		NotFound(c, "game not found")
		return
	}

	round, err := h.queries.GetRoundByID(ctx, roundID)
	if err != nil {
		NotFound(c, "round not found")
		return
	}

	if round.QuestionID.Valid {
		BadRequest(c, "question already set")
		return
	}

	// 驗證題目存在
	question, err := h.queries.GetQuestionByID(ctx, req.QuestionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			BadRequest(c, "question not found")
		} else {
			InternalServerError(c, "get question failed")
		}
		return
	}

	// 更新 round 的 question_id
	err = h.queries.SetRoundQuestionID(ctx, database.SetRoundQuestionIDParams{
		ID:         roundID,
		QuestionID: pgtype.Int8{Int64: req.QuestionID, Valid: true},
	})
	if err != nil {
		InternalServerError(c, "update round failed")
		return
	}

	h.hub.BroadcastToGame(game.Code, ws.WebSocketMessage{
		Type: "answer_time",
	})

	// 私訊題目內容給回答者
	h.hub.SendToPlayer(game.Code, round.AnswerPlayerID.Int64, ws.WebSocketMessage{
		Type: "round_question",
		Data: gin.H{
			"question": question,
		},
	})

	Success(c, gin.H{"message": "question set"})
}

type SubmitAnswerRequest struct {
	Answer string `json:"answer" binding:"required"`
}

func (h *RoundsHandler) SubmitAnswer(c *gin.Context) {
	ctx := c.Request.Context()
	code := c.Param("code")
	roundIDStr := c.Param("id")
	roundID, err := strconv.ParseInt(roundIDStr, 10, 64)
	if err != nil {
		BadRequest(c, "invalid round id")
		return
	}

	game, err := h.queries.GetGameByCode(ctx, code)
	if err != nil {
		NotFound(c, "game not found")
		return
	}

	var req SubmitAnswerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "invalid answer")
		return
	}

	// 確認 round 是否存在
	round, err := h.queries.GetRoundByID(ctx, roundID)
	if err != nil {
		NotFound(c, "round not found")
		return
	}

	// 更新 answer 欄位
	err = h.queries.SetRoundAnswer(ctx, database.SetRoundAnswerParams{
		ID:     round.ID,
		Answer: pgtype.Text{String: req.Answer, Valid: true},
	})
	if err != nil {
		InternalServerError(c, "update answer failed")
		return
	}

	// question, err := h.queries.GetQuestionByID(ctx, round.QuestionID.Int64)
	// if err != nil {
	// 	if errors.Is(err, sql.ErrNoRows) {
	// 		NotFound(c, "question not found")
	// 	} else {
	// 		InternalServerError(c, "failed to load question")
	// 	}
	// 	return
	// }

	// 廣播題目與答案
	h.hub.BroadcastToGame(game.Code, ws.WebSocketMessage{
		Type: "answer_submitted",
		Data: gin.H{
			"roundId": roundID,
			"answer":  req.Answer,
			// "question": question,
		},
	})

	Success(c, gin.H{"message": "answer submitted"})
}

func (h *RoundsHandler) DrawCard(c *gin.Context) {
	ctx := c.Request.Context()
	code := c.Param("code")
	roundID, err := utils.ParseID(c.Param("id"))
	if err != nil {
		BadRequest(c, "invalid round id")
		return
	}

	var req struct {
		Index    int   `json:"index"`
		PlayerID int64 `json:"playerId"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "invalid input")
		return
	}

	// 取得 game + round
	game, err := h.queries.GetGameByCode(ctx, code)
	if err != nil {
		NotFound(c, "game not found")
		return
	}

	round, err := h.queries.GetRoundByID(ctx, roundID)
	if err != nil {
		NotFound(c, "round not found")
		return
	}

	// ✅ 驗證是該輪的回答者
	if !round.AnswerPlayerID.Valid || round.AnswerPlayerID.Int64 != req.PlayerID {
		Forbidden(c, "not your turn to draw")
		return
	}

	// ✅ 驗證 index 合法
	if req.Index < 0 || req.Index >= len(round.Deck) {
		BadRequest(c, "invalid card index")
		return
	}

	// ✅ 判斷是否中 joker
	isJoker := round.Deck[req.Index] == "joker"

	// ✅ 更新 round 狀態為 done + 是否 joker
	err = h.queries.CompleteRound(ctx, database.CompleteRoundParams{
		ID:      round.ID,
		Status:  "revealed",
		IsJoker: pgtype.Bool{Bool: isJoker, Valid: true},
	})
	if err != nil {
		InternalServerError(c, "update round failed")
		return
	}

	// ✅ 撈出題目內容
	var questionContent string
	if round.QuestionID.Valid {
		q, err := h.queries.GetQuestionByID(ctx, round.QuestionID.Int64)
		if err == nil {
			questionContent = q
		}
	}

	// ✅ 廣播
	if isJoker {
		h.hub.BroadcastToGame(game.Code, ws.WebSocketMessage{
			Type: "joker_revealed",
			Data: gin.H{
				"roundId":  round.ID,
				"playerId": req.PlayerID,
				"question": questionContent,
				"answer":   round.Answer.String,
			},
		})
	} else {
		h.hub.BroadcastToGame(game.Code, ws.WebSocketMessage{
			Type: "player_safe",
			Data: gin.H{
				"roundId":  round.ID,
				"playerId": req.PlayerID,
			},
		})
	}

	// ✅ 回傳結果
	Success(c, gin.H{
		"isJoker": isJoker,
	})
}

func (h *RoundsHandler) CreateNextRound(c *gin.Context) {
	ctx := c.Request.Context()
	code := c.Param("code")

	var req struct {
		HostID         int64 `json:"hostId"`
		CurrentRoundID int64 `json:"currentRoundId"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "invalid input")
		return
	}

	// 取得遊戲
	game, err := h.queries.GetGameByCode(ctx, code)
	if err != nil {
		NotFound(c, "game not found")
		return
	}

	// 驗證 host 身分
	host, err := h.queries.GetPlayerByID(ctx, req.HostID)
	if err != nil || host.GameID != game.ID || !host.IsHost.Bool {
		Forbidden(c, "not a valid host")
		return
	}

	// 取得當前 round
	prevRound, err := h.queries.GetRoundByID(ctx, req.CurrentRoundID)
	if err != nil {
		NotFound(c, "previous round not found")
		return
	}

	// 取得該房所有玩家，按 joined_at 排序
	players, err := h.queries.ListPlayersByGameCode(ctx, game.Code)
	if err != nil || len(players) < 2 {
		InternalServerError(c, "not enough players")
		return
	}

	// 找出下一位回答者
	var nextPlayerID int64
	for i, p := range players {
		if p.ID == prevRound.AnswerPlayerID.Int64 {
			nextPlayerID = players[(i+1)%len(players)].ID
			break
		}
	}

	// 建立新 round
	deck := utils.CreateShuffledDeck(1, 3)
	round, err := h.queries.CreateRound(ctx, database.CreateRoundParams{
		GameID:           game.ID,
		QuestionPlayerID: prevRound.AnswerPlayerID.Int64,
		AnswerPlayerID:   pgtype.Int8{Int64: nextPlayerID, Valid: true},
		Deck:             deck,
	})
	if err != nil {
		InternalServerError(c, "failed to create next round")
		return
	}

	// 標記前一輪為 done
	_ = h.queries.FinishRound(ctx, req.CurrentRoundID)

	// 廣播下一輪開始
	h.hub.BroadcastToGame(game.Code, ws.WebSocketMessage{
		Type: "round_started",
		Data: gin.H{
			"roundId":      round.ID,
			"questionerId": round.QuestionPlayerID,
			"answererId":   round.AnswerPlayerID.Int64,
		},
	})

	Success(c, gin.H{
		"roundId":      round.ID,
		"questionerId": round.QuestionPlayerID,
		"answererId":   round.AnswerPlayerID.Int64,
	})
}

func (h *GamesHandler) GetQuestions(c *gin.Context) {
	ctx := c.Request.Context()
	code := c.Param("code")

	// 預設 limit 為 3 題
	limitStr := c.DefaultQuery("limit", "3")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 3
	}

	game, err := h.queries.GetGameByCode(ctx, code)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			NotFound(c, "game not found")
		} else {
			InternalServerError(c, "db error")
		}
		return
	}

	questions, err := h.queries.GetQuestionsByLevel(ctx, database.GetQuestionsByLevelParams{
		Level: game.Level,
		Limit: int32(limit),
	})

	if err != nil {
		InternalServerError(c, "failed to get questions")
		return
	}

	Success(c, transformToQuestionResponse(questions))
}

type QuestionResponse struct {
	ID      int64  `json:"id"`
	Content string `json:"content"`
}

func transformToQuestionResponse(questions []database.GetQuestionsByLevelRow) []QuestionResponse {
	var questionResponses []QuestionResponse
	for _, q := range questions {
		questionResponses = append(questionResponses, QuestionResponse{
			ID:      q.ID,
			Content: q.Content,
		})
	}
	return questionResponses
}
