package ws

import (
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	ID           int64           // 玩家 ID，供單播使用
	conn         *websocket.Conn // WebSocket 實際連線
	send         chan []byte     // 發送訊息用的 channel
	room         *Room           // 所屬房間
	OnDisconnect func(playerID int64)
}

func (c *Client) readPump() {
	defer c.disconnect()
	for {
		_, _, err := c.conn.ReadMessage()
		if err != nil {
			log.Printf("read error: %v", err)
			break
		}
	}
}

func (c *Client) writePump() {
	for msg := range c.send {
		err := c.conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Printf("write error: %v", err)
			break
		}
	}
}

func (c *Client) disconnect() {
	c.room.leave <- c

	_ = c.conn.Close()
	if c.OnDisconnect != nil {
		c.OnDisconnect(c.ID)
	}

}
