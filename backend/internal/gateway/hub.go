package gateway

import (
	"gopherspace/internal/core_ai"
)

type Hub struct {
	Clients    map[*Client]bool // 账本：记录当前所有在线的 Client
	Broadcast  chan []byte      // 广播队列：塞进这里的消息会发给所有人
	Register   chan *Client     // 注册队列：新学生连进来了，往这里丢
	Unregister chan *Client     // 注销队列：学生网页关闭了，往这里丢
	AiModel    *core_ai.CoreAiService
}

func NewHub(ai *core_ai.CoreAiService) *Hub {
	return &Hub{
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		AiModel:    ai,
	}
}
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client] = true
		case client := <-h.Unregister:
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.Send)
			}
		case message := <-h.Broadcast:
			for client := range h.Clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.Clients, client)
				}
			}
		}

	}
}
