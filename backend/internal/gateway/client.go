package gateway

import (
	"context"
	"encoding/json"
	"io"
	"log"

	"gopherspace/internal/core_ai"

	"github.com/gorilla/websocket"
)

type Client struct {
	Hub       *Hub
	Conn      *websocket.Conn
	Send      chan []byte
	ChatGraph *core_ai.ChatGraph
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
			stream, err := c.ChatGraph.Stream(context.Background(), core_ai.ChatRequest{
				UserMessage: wsReq.Content,
				History:     nil,
			})
			if err != nil {
				log.Println("ChatGraph.Stream 启动失败:", err)
				c.Send <- toJSON(WbResponse{Type: "ai", Content: "抱歉，AI 服务暂时不可用，请稍后再试。", Done: true})
				continue
			}
			hasContent := false
			for {
				chunk, err := stream.Recv()
				if err != nil {
					if err != io.EOF {
						log.Printf("stream.Recv 错误: %v", err)
					}
					break
				}
				log.Printf("收到 chunk: %q", chunk.Content)
				if chunk != nil && chunk.Content != "" {
					hasContent = true
					c.Send <- toJSON(WbResponse{Type: "ai", Content: chunk.Content, Done: false})
				}
			}
			stream.Close()
			log.Printf("流式结束, hasContent=%v", hasContent)
			if !hasContent {
				// 先发内容（Done: false），再发结束标记
				c.Send <- toJSON(WbResponse{Type: "ai", Content: "抱歉，AI 未能生成回复，请换个方式描述你的问题。", Done: false})
			}
			c.Send <- toJSON(WbResponse{Type: "ai", Content: "", Done: true})
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
