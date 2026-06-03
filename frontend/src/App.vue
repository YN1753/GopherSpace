<script setup>
import { ref, reactive, onMounted, onUnmounted, nextTick, computed } from 'vue'
import ChatMessage from './components/ChatMessage.vue'

const messages = reactive([])
const inputText = ref('')
const isConnected = ref(false)
const isStreaming = ref(false)
const messagesContainer = ref(null)
let ws = null
let currentAiMessage = null
let reconnectTimer = null
let wasConnected = false

// Connection status
const statusText = computed(() => {
  if (isConnected.value) return '已连接'
  return '未连接'
})

function connect() {
  const wsUrl = 'ws://localhost:8181/v1/ws'
  ws = new WebSocket(wsUrl)

  ws.onopen = () => {
    isConnected.value = true
    wasConnected = true
  }

  ws.onmessage = (event) => {
    try {
      const data = JSON.parse(event.data)
      if (data.type === 'ai') {
        if (data.done) {
          // AI response finished
          isStreaming.value = false
          currentAiMessage = null
        } else {
          // Streaming chunk
          if (data.content) {
            if (!currentAiMessage) {
              currentAiMessage = {
                id: Date.now(),
                role: 'ai',
                content: data.content.replace(/^\n+/, ''),
                timestamp: new Date(),
              }
              messages.push(currentAiMessage)
            } else {
              currentAiMessage.content += data.content
            }
          }
          scrollToBottom()
        }
      }
    } catch (e) {
      // Plain text message
      addSystemMessage(event.data)
    }
  }

  ws.onclose = () => {
    isConnected.value = false
    isStreaming.value = false
    if (wasConnected) {
      wasConnected = false
      addSystemMessage('连接已断开，正在重连...')
    }
    clearTimeout(reconnectTimer)
    reconnectTimer = setTimeout(connect, 3000)
  }

  ws.onerror = () => {
    isConnected.value = false
  }
}

function sendMessage() {
  const text = inputText.value.trim()
  if (!text || !isConnected.value || isStreaming.value) return

  // Add user message
  messages.push({
    id: Date.now(),
    role: 'user',
    content: text,
    timestamp: new Date(),
  })

  // Send to backend
  ws.send(JSON.stringify({
    type: 'ai',
    content: text,
  }))

  inputText.value = ''
  isStreaming.value = true
  scrollToBottom()
}

function addSystemMessage(text) {
  messages.push({
    id: Date.now() + Math.random(),
    role: 'system',
    content: text,
    timestamp: new Date(),
  })
  scrollToBottom()
}

function scrollToBottom() {
  nextTick(() => {
    if (messagesContainer.value) {
      messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
    }
  })
}

function handleKeydown(e) {
  if (e.key === 'Enter' && !e.shiftKey) {
    e.preventDefault()
    sendMessage()
  }
}

onMounted(() => {
  connect()
})

onUnmounted(() => {
  clearTimeout(reconnectTimer)
  if (ws) {
    ws.onclose = null
    ws.close()
  }
})
</script>

<template>
  <div class="app">
    <!-- Header -->
    <header class="header">
      <div class="header-content">
        <div class="logo-section">
          <div class="logo">
            <svg width="28" height="28" viewBox="0 0 28 28" fill="none">
              <circle cx="14" cy="14" r="12" stroke="url(#grad)" stroke-width="2.5" fill="none"/>
              <circle cx="14" cy="14" r="5" fill="url(#grad)"/>
              <defs>
                <linearGradient id="grad" x1="0" y1="0" x2="28" y2="28">
                  <stop offset="0%" stop-color="#6c5ce7"/>
                  <stop offset="100%" stop-color="#a29bfe"/>
                </linearGradient>
              </defs>
            </svg>
          </div>
          <div class="logo-text">
            <h1>GopherSpace</h1>
            <span class="subtitle">AI Assistant</span>
          </div>
        </div>
        <div class="status-section">
          <span class="status-dot" :class="{ connected: isConnected }"></span>
          <span class="status-text">{{ statusText }}</span>
        </div>
      </div>
    </header>

    <!-- Messages -->
    <main class="messages-area" ref="messagesContainer">
      <div class="messages-inner">
        <!-- Welcome -->
        <div v-if="messages.length === 0" class="welcome">
          <div class="welcome-icon">
            <svg width="48" height="48" viewBox="0 0 48 48" fill="none">
              <circle cx="24" cy="24" r="20" stroke="url(#welcomeGrad)" stroke-width="2" fill="none" opacity="0.5"/>
              <circle cx="24" cy="24" r="8" fill="url(#welcomeGrad)"/>
              <defs>
                <linearGradient id="welcomeGrad" x1="0" y1="0" x2="48" y2="48">
                  <stop offset="0%" stop-color="#6c5ce7"/>
                  <stop offset="100%" stop-color="#a29bfe"/>
                </linearGradient>
              </defs>
            </svg>
          </div>
          <h2>开始对话</h2>
          <p>输入任何问题，AI 将为你解答</p>
        </div>

        <!-- Message List -->
        <ChatMessage
          v-for="msg in messages"
          :key="msg.id"
          :message="msg"
        />

        <!-- Streaming indicator -->
        <div v-if="isStreaming && !currentAiMessage" class="typing-indicator">
          <div class="typing-dots">
            <span></span><span></span><span></span>
          </div>
        </div>
      </div>
    </main>

    <!-- Input -->
    <footer class="input-area">
      <div class="input-wrapper">
        <textarea
          v-model="inputText"
          @keydown="handleKeydown"
          :placeholder="isConnected ? '输入消息...' : '连接中...'"
          :disabled="!isConnected || isStreaming"
          rows="1"
          ref="inputRef"
        ></textarea>
        <button
          class="send-btn"
          @click="sendMessage"
          :disabled="!inputText.trim() || !isConnected || isStreaming"
        >
          <svg width="20" height="20" viewBox="0 0 20 20" fill="none">
            <path d="M3 10L17 3L10 17L9 11L3 10Z" fill="currentColor"/>
          </svg>
        </button>
      </div>
      <div class="input-hint">
        <span v-if="isStreaming" class="streaming-hint">
          <span class="pulse"></span> AI 正在回复...
        </span>
        <span v-else>Enter 发送 · Shift+Enter 换行</span>
      </div>
    </footer>
  </div>
</template>

<style scoped>
.app {
  height: 100vh;
  display: flex;
  flex-direction: column;
  background: var(--bg-primary);
}

/* Header */
.header {
  border-bottom: 1px solid var(--border);
  background: var(--bg-secondary);
  backdrop-filter: blur(20px);
  position: sticky;
  top: 0;
  z-index: 10;
}

.header-content {
  max-width: 800px;
  margin: 0 auto;
  padding: 16px 24px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.logo-section {
  display: flex;
  align-items: center;
  gap: 12px;
}

.logo {
  display: flex;
  align-items: center;
}

.logo-text h1 {
  font-size: 18px;
  font-weight: 600;
  color: var(--text-primary);
  letter-spacing: -0.3px;
}

.logo-text .subtitle {
  font-size: 12px;
  color: var(--text-muted);
  font-weight: 400;
}

.status-section {
  display: flex;
  align-items: center;
  gap: 8px;
}

.status-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: var(--error);
  transition: background 0.3s;
}

.status-dot.connected {
  background: var(--success);
  box-shadow: 0 0 8px rgba(0, 206, 201, 0.4);
}

.status-text {
  font-size: 13px;
  color: var(--text-muted);
}

/* Messages */
.messages-area {
  flex: 1;
  overflow-y: auto;
  padding: 24px;
}

.messages-inner {
  max-width: 800px;
  margin: 0 auto;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.welcome {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 80px 20px;
  text-align: center;
}

.welcome-icon {
  margin-bottom: 24px;
  opacity: 0.8;
}

.welcome h2 {
  font-size: 24px;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 8px;
}

.welcome p {
  font-size: 15px;
  color: var(--text-muted);
}

/* Typing indicator */
.typing-indicator {
  display: flex;
  padding: 16px 20px;
}

.typing-dots {
  display: flex;
  gap: 6px;
  align-items: center;
}

.typing-dots span {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: var(--accent);
  opacity: 0.4;
  animation: typingBounce 1.4s ease-in-out infinite;
}

.typing-dots span:nth-child(2) {
  animation-delay: 0.2s;
}

.typing-dots span:nth-child(3) {
  animation-delay: 0.4s;
}

@keyframes typingBounce {
  0%, 80%, 100% {
    transform: scale(0.8);
    opacity: 0.4;
  }
  40% {
    transform: scale(1.2);
    opacity: 1;
  }
}

/* Input */
.input-area {
  border-top: 1px solid var(--border);
  background: var(--bg-secondary);
  padding: 16px 24px 20px;
}

.input-wrapper {
  max-width: 800px;
  margin: 0 auto;
  display: flex;
  align-items: flex-end;
  gap: 12px;
  background: var(--bg-input);
  border: 1px solid var(--border);
  border-radius: var(--radius-lg);
  padding: 8px 8px 8px 20px;
  transition: border-color 0.2s, box-shadow 0.2s;
}

.input-wrapper:focus-within {
  border-color: var(--accent);
  box-shadow: 0 0 0 3px var(--accent-glow);
}

textarea {
  flex: 1;
  background: transparent;
  border: none;
  outline: none;
  color: var(--text-primary);
  font-family: var(--font);
  font-size: 15px;
  line-height: 1.5;
  resize: none;
  max-height: 120px;
  padding: 8px 0;
}

textarea::placeholder {
  color: var(--text-muted);
}

textarea:disabled {
  opacity: 0.5;
}

.send-btn {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  border: none;
  background: var(--accent);
  color: white;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
  flex-shrink: 0;
}

.send-btn:hover:not(:disabled) {
  background: var(--accent-light);
  transform: scale(1.05);
}

.send-btn:disabled {
  opacity: 0.3;
  cursor: not-allowed;
}

.input-hint {
  max-width: 800px;
  margin: 8px auto 0;
  font-size: 12px;
  color: var(--text-muted);
  text-align: center;
}

.streaming-hint {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  color: var(--accent-light);
}

.pulse {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: var(--accent);
  animation: pulse 1.5s ease-in-out infinite;
}

@keyframes pulse {
  0%, 100% { opacity: 0.4; transform: scale(0.8); }
  50% { opacity: 1; transform: scale(1.2); }
}

/* Responsive */
@media (max-width: 640px) {
  .header-content {
    padding: 12px 16px;
  }
  .messages-area {
    padding: 16px;
  }
  .input-area {
    padding: 12px 16px 16px;
  }
}
</style>
