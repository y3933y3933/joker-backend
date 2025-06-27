package api

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"math/rand"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/y3933y3933/joker/internal/database"
	"github.com/y3933y3933/joker/internal/utils"
	"github.com/y3933y3933/joker/internal/ws"
)

type GamesHandler struct {
	logger  *slog.Logger
	queries *database.Queries
	hub     *ws.Hub
}

func NewGamesHandler(queries *database.Queries, logger *slog.Logger, hub *ws.Hub) *GamesHandler {
	return &GamesHandler{
		logger:  logger,
		queries: queries,
		hub:     hub,
	}
}

type CreateGameRequest struct {
	Level string `json:"level" binding:"required,oneof=easy normal spicy"`
}

type CreateGameResponse struct {
	ID        int64     `json:"id"`
	Code      string    `json:"code"`
	Level     string    `json:"level"`
	CreatedAt time.Time `json:"createdAt"`
}

func (h *GamesHandler) CreateGame(c *gin.Context) {
	var req CreateGameRequest
	if err := bindCreateGameRequest(c, &req); err != nil {
		return
	}

	ctx := c.Request.Context()

	code, err := generateGameCode(ctx, h)
	if err != nil {
		handleGameCodeError(c, h.logger, err)
		return
	}

	game, err := h.queries.CreateGame(ctx, database.CreateGameParams{
		Code:   code,
		Level:  req.Level,
		Status: "waiting",
	})
	if err != nil {
		h.logger.Error("create game fail: ", err)
		InternalServerError(c, "failed to create game")
		return
	}

	Success(c, CreateGameResponse{
		ID:        game.ID,
		Code:      game.Code,
		Level:     game.Level,
		CreatedAt: game.CreatedAt.Time,
	})

}

func generateGameCode(ctx context.Context, h *GamesHandler) (string, error) {
	return utils.GenerateUniqueGameCode(ctx, h.queries, 6, 5)
}

func bindCreateGameRequest(c *gin.Context, req *CreateGameRequest) error {
	if err := c.ShouldBindJSON(req); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			FailedValidation(c, ve.Error())
		} else {
			BadRequest(c, "Invalid request format")
		}
		return err
	}
	return nil
}

func handleGameCodeError(c *gin.Context, logger *slog.Logger, err error) {
	logger.Error("generate unique game code error: ", err)
	if errors.Is(err, utils.ErrGenerateCode) {
		InternalServerError(c, "code collision, try again")
	} else {
		InternalServerError(c, "DB error")
	}
}

type GameResponse struct {
	ID        int64     `json:"id"`
	Code      string    `json:"code"`
	Level     string    `json:"level"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
}

func (h *GamesHandler) GetGame(c *gin.Context) {
	ctx := c.Request.Context()
	code := c.Param("code")

	game, err := h.queries.GetGameByCode(ctx, code)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			NotFound(c, "Game not found")
		} else {
			InternalServerError(c, "DB error")
		}
		return
	}

	Success(c, GameResponse{
		ID:        game.ID,
		Code:      game.Code,
		Level:     game.Level,
		Status:    game.Status,
		CreatedAt: game.CreatedAt.Time,
	})
}

func (h *GamesHandler) EndGame(c *gin.Context) {
	ctx := c.Request.Context()
	gameCode := c.Param("code")
	playerIDStr := c.GetHeader("X-Player-ID")
	if playerIDStr == "" {
		BadRequest(c, "missing X-Player-ID header")
		return
	}
	playerID, err := utils.ParseID(playerIDStr)
	if err != nil {
		BadRequest(c, "invalid player ID")
		return
	}

	player, err := h.queries.GetPlayerByGameCodeAndID(c, database.GetPlayerByGameCodeAndIDParams{
		Code: gameCode,
		ID:   playerID,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			NotFound(c, "player not found in game")
		} else {
			h.logger.Error("failed to get player", slog.Any("err", err))
			InternalServerError(c, "failed to verify player")
		}
		return
	}
	if !player.IsHost.Valid || (player.IsHost.Valid && !player.IsHost.Bool) {
		Forbidden(c, "only host can end the game")
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

	err = h.queries.UpdateGameStatus(ctx, database.UpdateGameStatusParams{
		ID:     game.ID,
		Status: "ended",
	})
	if err != nil {
		InternalServerError(c, "failed to end game")
		return
	}

	// 廣播遊戲結束
	h.hub.BroadcastToGame(game.Code, ws.WebSocketMessage{
		Type: "game_ended",
		Data: gin.H{
			"gameId":   game.ID,
			"gameCode": game.Code,
		},
	})

	c.Status(http.StatusOK)
}

type GameSummaryResponse struct {
	Game       GameInfo        `json:"game"`
	Players    []PlayerSummary `json:"players"`
	Rounds     []RoundSummary  `json:"rounds"`
	Highlights Highlights      `json:"highlights"`
}

type GameInfo struct {
	Code            string `json:"code"`
	Level           string `json:"level"`
	TotalRounds     int    `json:"totalRounds"`
	TotalJokerCards int    `json:"totalJokerCards"`
}

type PlayerSummary struct {
	PlayerID          int64  `json:"playerId"`
	Nickname          string `json:"nickname"`
	QuestionsAsked    int64  `json:"questionsAsked"`
	QuestionsAnswered int64  `json:"questionsAnswered"`
	JokerCardsDrawn   int64  `json:"jokerCardsDrawn"`
}

type RoundSummary struct {
	RoundID                int64  `json:"roundId"`
	Question               string `json:"question"`
	Answer                 string `json:"answer"`
	IsJoker                bool   `json:"isJoker"`
	QuestionPlayerID       int64  `json:"questionPlayerId"`
	QuestionPlayerNickname string `json:"questionPlayerNickname"`
	AnswerPlayerID         int64  `json:"answerPlayerId"`
	AnswerPlayerNickname   string `json:"answerPlayerNickname"`
}

type Highlights struct {
	LuckiestPlayer   PlayerSummary `json:"luckiestPlayer"`
	UnluckiestPlayer PlayerSummary `json:"unluckiestPlayer"`
	BestQAs          []QAItem      `json:"bestQAs"`
}

type QAItem struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

func (h *GamesHandler) GetGameSummary(c *gin.Context) {
	code := c.Param("code")

	game, err := h.queries.GetGameByCode(c, code)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			NotFound(c, "Game not found or not ended")
		} else {
			InternalServerError(c, "Failed to get game")
		}
		return
	}

	roundsRaw, err := h.queries.GetRoundSummaryByGameID(c, game.ID)
	if err != nil {
		InternalServerError(c, "Failed to get rounds")
		return
	}

	playersRaw, err := h.queries.GetPlayerStatsByGameID(c, game.ID)
	if err != nil {
		InternalServerError(c, "Failed to get player stats")
		return
	}

	// 轉換玩家資料 & 統計鬼牌
	var (
		playerList       []PlayerSummary
		luckiestPlayer   PlayerSummary
		unluckiestPlayer PlayerSummary
		totalJoker       int
		minJoker         = int(^uint(0) >> 1) // 最大整數
		maxJoker         = -1
	)

	for _, p := range playersRaw {
		player := PlayerSummary{
			PlayerID:          p.PlayerID,
			Nickname:          p.Nickname,
			QuestionsAsked:    p.QuestionsAsked,
			QuestionsAnswered: p.QuestionsAnswered,
			JokerCardsDrawn:   p.JokerCardsDrawn,
		}
		playerList = append(playerList, player)

		joker := int(p.JokerCardsDrawn)
		totalJoker += joker

		if joker < minJoker {
			minJoker = joker
			luckiestPlayer = player
		}
		if joker > maxJoker {
			maxJoker = joker
			unluckiestPlayer = player
		}
	}

	// 轉換回合資料
	var (
		roundList   []RoundSummary
		eligibleQAs []QAItem
	)

	for _, r := range roundsRaw {
		isJoker := r.IsJoker.Valid && r.IsJoker.Bool
		answer := ""
		if r.Answer.Valid {
			answer = r.Answer.String
		}

		round := RoundSummary{
			RoundID:                r.RoundID,
			Question:               r.Question,
			Answer:                 answer,
			IsJoker:                isJoker,
			QuestionPlayerID:       r.QuestionPlayerID,
			QuestionPlayerNickname: r.QuestionPlayerNickname,
			AnswerPlayerID:         r.AnswerPlayerID,
			AnswerPlayerNickname:   r.AnswerPlayerNickname,
		}
		roundList = append(roundList, round)

		if isJoker && r.Answer.Valid {
			eligibleQAs = append(eligibleQAs, QAItem{
				Question: r.Question,
				Answer:   r.Answer.String,
			})
		}
	}

	// 隨機選 3 筆 QA
	rand.Shuffle(len(eligibleQAs), func(i, j int) {
		eligibleQAs[i], eligibleQAs[j] = eligibleQAs[j], eligibleQAs[i]
	})
	if len(eligibleQAs) > 3 {
		eligibleQAs = eligibleQAs[:3]
	}

	// 組裝回傳
	response := GameSummaryResponse{
		Game: GameInfo{
			Code:            game.Code,
			Level:           game.Level,
			TotalRounds:     len(roundList),
			TotalJokerCards: totalJoker,
		},
		Players: playerList,
		Rounds:  roundList,
		Highlights: Highlights{
			LuckiestPlayer:   luckiestPlayer,
			UnluckiestPlayer: unluckiestPlayer,
			BestQAs:          eligibleQAs,
		},
	}

	Success(c, response)
	// c.JSON(http.StatusOK, response)
}
