# GopherSpace

基于 Go + Vue 3 的 AI 对话应用，支持 WebSocket 流式输出。

## 技术栈

**后端**

- Go + Gin（HTTP 框架）
- gorilla/websocket（WebSocket 通信）
- Eino 框架 + Claude 兼容组件（对接 MiniMax 大模型）

**前端**

- Vue 3 + Vite
- 原生 WebSocket API
- 深色主题 UI

## 项目结构

```
gopherspace/
├── README.md
├── backend/
│   ├── cmd/
│   │   └── main.go                 # 入口
│   ├── config/                     # 配置（待实现）
│   ├── internal/
│   │   ├── core_ai/
│   │   │   └── chat.go             # AI 服务层
│   │   └── gateway/
│   │       ├── client.go           # WebSocket 客户端（读写泵）
│   │       ├── hub.go              # 连接管理中心
│   │       ├── types.go            # 消息协议定义
│   │       ├── ws_client.go        # 预留
│   │       └── ws_handler.go       # HTTP → WebSocket 升级
│   ├── go.mod
│   ├── go.sum
│   └── minimax.go                  # MiniMax 模型初始化
└── frontend/
    ├── index.html
    ├── vite.config.js
    └── src/
        ├── main.js
        ├── App.vue                 # 主组件（WebSocket + 聊天逻辑）
        ├── style.css               # 全局样式
        └── components/
            └── ChatMessage.vue     # 消息气泡组件
```

## 消息协议

前后端通过 JSON 格式通信：

**客户端 → 服务端**

```json
{"type": "ai", "content": "你好"}
```

**服务端 → 客户端（流式）**

```json
{"type": "ai", "content": "你好！", "done": false}
{"type": "ai", "content": "有什么可以帮你的？", "done": false}
{"type": "ai", "content": "", "done": true}
```

- `done: false` — AI 正在输出，`content` 为当前 chunk
- `done: true` — AI 回复结束

## 快速启动

**后端**

```bash
cd backend
go run cmd/main.go
```

默认监听 `:8181`，需要配置环境变量 `MINIMAX_API_KEY`。

**前端**

```bash
cd frontend
npm install
npm run dev
```

浏览器打开 `http://localhost:3000`。

## 当前进度

- [x] WebSocket 网关（Hub 模式连接管理）
- [x] MiniMax 大模型流式对接
- [x] 前端深色主题聊天 UI
- [x] AI 流式逐字输出
- [x] 消息类型协议（ai / broadcast）
- [ ] 前端广播消息支持
- [ ] 会话上下文记忆（多轮对话）
- [ ] 用户认证
- [ ] 配置管理（config 文件）
- [ ] 数据持久化
