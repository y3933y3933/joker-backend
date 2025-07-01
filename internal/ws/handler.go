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
}

// NewHandler 用來建立新的 WebSocket handler
func NewHandler(hub *Hub, logger *slog.Logger, playerService *service.PlayerService) *Handler {
	return &Handler{Hub: hub, Logger: logger, PlayerService: playerService}

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
		OnDisconnect: func(gameCode string, playerID int64) {
			left, newHost, err := h.PlayerService.LeaveGame(context.Background(), playerID)
			if err != nil {
				h.Logger.Error("Failed to remove player", "error", err)
			}

			msg1, _ := NewWSMessage(MsgPlayerLeft, PlayerLeftPayload{
				ID:       left.ID,
				Nickname: left.Nickname,
			})
			room.Broadcast(msg1)

			if newHost != nil {
				msg2, _ := NewWSMessage(MsgHostTransferred, HostTransferredPayload{
					ID:       newHost.ID,
					Nickname: newHost.Nickname,
				})
				room.Broadcast(msg2)
			}
		},
	}

	room.join <- client

	go client.writePump()
	go client.readPump()
}
