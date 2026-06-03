package gateway

import (
	"context"
	"encoding/json"
	"gopherspace/internal/core_ai"
	"log"

	"github.com/cloudwego/eino/schema"
	"github.com/gorilla/websocket"
)

type Client struct {
	Hub     *Hub
	Conn    *websocket.Conn
	Send    chan []byte
	AiModel *core_ai.CoreAiService
}

func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister <- c
		_ = c.Conn.Close()
	}()
	for {
		_, message, err := c.Conn.ReadMessage()
		var wsReq WbRequest
		if err != nil {
			log.Println("读取失败:", err)
			break
		}
		err = json.Unmarshal(message, &wsReq)
		if err != nil {
			log.Println("转换格式失败:", err)
		}
		switch wsReq.Type {
		case "ai":
			messages := []*schema.Message{
				schema.SystemMessage("You are a helpful assistant."),
				schema.UserMessage(wsReq.Content),
			}
			stream, err := c.AiModel.Model.Stream(context.Background(), messages)
			if err != nil {
				log.Println(err)
				continue
			}
			for {
				chunk, err := stream.Recv()
				if err != nil {
					break
				}
				c.Send <- toJSON(WbResponse{Type: "ai", Content: chunk.Content, Done: false})
			}
			c.Send <- toJSON(WbResponse{Type: "ai", Content: "", Done: true})
			stream.Close()
		case "broadcast":
			c.Hub.Broadcast <- []byte(wsReq.Content)

		}
	}
}
func (c *Client) WritePump() {
	defer func() {
		_ = c.Conn.Close()
	}()
	for {
		message, ok := <-c.Send
		if !ok {
			_ = c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
			return
		}
		err := c.Conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Println(err)
			return
		}
	}
}
