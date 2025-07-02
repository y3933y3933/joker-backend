package ws

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/y3933y3933/joker/internal/service"
	"github.com/y3933y3933/joker/internal/store"
	"github.com/y3933y3933/joker/internal/utils/errx"
	"github.com/y3933y3933/joker/internal/utils/httpx"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// 開放所有來源，正式環境建議加驗證
		return true
	},
}

// Handler struct 用來包裝 Hub 實例
type Handler struct {
	Hub           *Hub
	Logger        *slog.Logger
	PlayerService *service.PlayerService
	GameService   *service.GameService
	RoundService  *service.RoundService
}

// NewHandler 用來建立新的 WebSocket handler
func NewHandler(hub *Hub, logger *slog.Logger, playerService *service.PlayerService, gameService *service.GameService, roundService *service.RoundService) *Handler {
	return &Handler{Hub: hub, Logger: logger, PlayerService: playerService, GameService: gameService, RoundService: roundService}

}

func (h *Handler) ServeWS(c *gin.Context) {
	gameCode := c.Param("code")
	playerIDStr := c.Query("player_id")
	playerID, err := strconv.ParseInt(playerIDStr, 10, 64)
	if err != nil {
		httpx.ServerErrorResponse(c, h.Logger, errors.New("invalid player id"))
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		h.Logger.Error("ws upgrade error: ", err)
		return
	}
	fmt.Printf("ServeWS: playerID=%d (query=%s)\n", playerID, playerIDStr)

	room := h.Hub.GetRoom(gameCode)
	if room == nil {
		room = h.Hub.CreateRoom(gameCode)
	}
	fmt.Printf("ServeWS: room=%p\n", room)

	client := &Client{
		ID:   playerID,
		conn: conn,
		send: make(chan []byte, 256),
		room: room,
		OnDisconnect: func(playerID int64) {
			ctx := context.Background()
			left, newHost, err := h.PlayerService.LeaveGame(ctx, playerID)
			// success
			if err == nil {
				// ✅ 廣播離開訊息
				msg1, _ := NewWSMessage(MsgPlayerLeft, PlayerLeftPayload{
					ID:       left.ID,
					Nickname: left.Nickname,
				})
				room.Broadcast(msg1)

				// ✅ 如果有 host 轉移，廣播
				if newHost != nil {
					msg2, _ := NewWSMessage(MsgHostTransferred, HostTransferredPayload{
						ID:       newHost.ID,
						Nickname: newHost.Nickname,
					})
					room.Broadcast(msg2)
				}

				// // 如果房內沒人，自動刪除 game
				// if room.PlayerCount() == 0 {
				// 	h.GameService.DeleteGameIfEmpty(ctx, gameCode)
				// 	h.Hub.DeleteRoom(gameCode)
				// }
			} else {
				if errors.Is(err, errx.ErrGameAlreadyStarted) {
					h.Logger.Info("Player disconnected during game", "playerID", playerID)
					// ✅ 遊戲中：標記為離線
					err = h.PlayerService.MarkPlayerDisconnected(ctx, playerID)
					if err != nil {
						h.Logger.Error("MarkPlayerDisconnected failed", "error", err)
						return
					}

					// 查找目前回合
					game, err := h.GameService.GetGameByCode(ctx, gameCode)
					if err != nil {
						h.Logger.Error("FindByID(game) failed", "error", err)
						return
					}

					round, err := h.RoundService.FindLastRoundByGameID(ctx, game.ID)
					if err != nil {
						h.Logger.Error("FindCurrentRoundByGameID failed", "error", err)
						return
					}

					isQuestionPlayer := round.Status == store.RoundStatusWaitingForQuestion && round.QuestionPlayerID == playerID
					isAnswerPlayer := (round.Status == store.RoundStatusWaitingForAnswer || round.Status == store.RoundStatusWaitingForDraw) && round.AnswerPlayerID == playerID

					if isQuestionPlayer || isAnswerPlayer {
						newRound, err := h.RoundService.SkipRound(ctx, game, round.ID)
						if err != nil {
							if errors.Is(err, errx.ErrNotEnoughPlayers) {
								// 結束遊戲
								err := h.GameService.EndGame(ctx, gameCode, store.GameStatusPlaying)
								if err != nil {
									h.Logger.Error("End game failed", "error", err)
									return
								}
								msg, _ := NewWSMessage(MsgTypeGameEnded, gin.H{"gameCode": game.Code})
								room.Broadcast(msg)

								return
							}
							h.Logger.Error("SkipRound failed", "error", err)
							return
						}

						msg, _ := NewWSMessage(MsgTypeRoundSkipped, RoundSkippedPayload{
							Reason: "disconnect",
							RoundStartedPayload: RoundStartedPayload{
								RoundID:          newRound.ID,
								QuestionPlayerID: newRound.QuestionPlayerID,
								AnswererID:       newRound.AnswerPlayerID,
							},
						})
						room.Broadcast(msg)

					}

				} else {
					h.Logger.Error("LeaveGame failed", "error", err)
				}
				return
			}

		},
	}

	room.join <- client

	go client.writePump()
	go client.readPump()
}
