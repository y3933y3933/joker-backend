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
}

// NewHandler 用來建立新的 WebSocket handler
func NewHandler(hub *Hub, logger *slog.Logger, playerService *service.PlayerService, gameService *service.GameService) *Handler {
	return &Handler{Hub: hub, Logger: logger, PlayerService: playerService, GameService: gameService}

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
			if err != nil {
				if errors.Is(err, errx.ErrGameAlreadyStarted) {
					// ✅ 遊戲中，不移除玩家，但仍可以標記離線（未來用）
					h.Logger.Info("Player disconnected during game", "playerID", playerID)
					// TODO: 可在這裡呼叫 MarkPlayerDisconnected()
				} else {
					h.Logger.Error("LeaveGame failed", "error", err)
				}
				return
			}

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

			// 如果房內沒人，自動刪除 game
			if room.PlayerCount() == 0 {
				h.GameService.DeleteGameIfEmpty(ctx, gameCode)
				h.Hub.DeleteRoom(gameCode)
			}
		},
	}

	room.join <- client

	go client.writePump()
	go client.readPump()
}

func (h *Handler) OnDisconnect(playerID int64) {}
