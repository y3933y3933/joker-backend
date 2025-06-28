package ws

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/y3933y3933/joker/internal/api"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// 開放所有來源，正式環境建議加驗證
		return true
	},
}

// Handler struct 用來包裝 Hub 實例
type Handler struct {
	Hub    *Hub
	Logger *slog.Logger
}

// NewHandler 用來建立新的 WebSocket handler
func NewHandler(hub *Hub, logger *slog.Logger) *Handler {
	return &Handler{Hub: hub, Logger: logger}

}

func (h *Handler) ServeWS(c *gin.Context) {
	gameCode := c.Param("code")
	playerIDStr := c.Query("player_id")
	playerID, err := strconv.ParseInt(playerIDStr, 10, 64)
	if err != nil {
		api.ServerErrorResponse(c, h.Logger, errors.New("invalid player id"))
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		h.Logger.Error("ws upgrade error: ", err)
		return
	}

	room := h.Hub.GetRoom(gameCode)
	if room == nil {
		room = h.Hub.CreateRoom(gameCode)
	}

	client := &Client{
		ID:   playerID,
		conn: conn,
		send: make(chan []byte, 256),
		room: room,
	}

	room.join <- client

	go client.writePump()
	go client.readPump()
}
