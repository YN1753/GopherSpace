package gateway

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true }, // 允许跨域
}

func HttpToWSEntry(hub *Hub) gin.HandlerFunc {
	return func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			return
		}

		client := &Client{
			Hub:       hub,
			Conn:      conn,
			Send:      make(chan []byte, 256),
			ChatGraph: hub.ChatGraph,
		}


		hub.Register <- client

		go client.WritePump()
		go client.ReadPump()
	}
}
