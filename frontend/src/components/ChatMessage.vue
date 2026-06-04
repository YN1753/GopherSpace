<script setup>
import { computed } from 'vue'
import { marked } from 'marked'
import DOMPurify from 'isomorphic-dompurify'

const props = defineProps({
  message: {
    type: Object,
    required: true,
  },
})

const isUser = computed(() => props.message.role === 'user')
const isSystem = computed(() => props.message.role === 'system')
const isAi = computed(() => props.message.role === 'ai')

const timeStr = computed(() => {
  const d = new Date(props.message.timestamp)
  return d.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })
})

// Parse markdown to HTML for AI messages
const renderedContent = computed(() => {
  if (isUser.value) {
    return props.message.content
  }
  const rawHtml = marked.parse(props.message.content || '', {
    breaks: true,
    gfm: true,
    async: false,
  })
  return DOMPurify.sanitize(rawHtml)
})
</script>

<template>
  <div
    class="message"
    :class="{
      'message-user': isUser,
      'message-ai': isAi,
      'message-system': isSystem,
    }"
  >
    <!-- System message -->
    <div v-if="isSystem" class="system-msg">
      <span>{{ message.content }}</span>
    </div>

    <!-- User / AI message -->
    <template v-else>
      <div class="avatar" :class="{ 'avatar-user': isUser, 'avatar-ai': isAi }">
        <template v-if="isUser">
          <svg width="18" height="18" viewBox="0 0 18 18" fill="none">
            <circle cx="9" cy="6" r="3.5" stroke="currentColor" stroke-width="1.5"/>
            <path d="M2.5 16.5C2.5 13 5.5 10.5 9 10.5C12.5 10.5 15.5 13 15.5 16.5" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
          </svg>
        </template>
        <template v-else>
          <svg width="18" height="18" viewBox="0 0 18 18" fill="none">
            <circle cx="9" cy="9" r="7" stroke="currentColor" stroke-width="1.5" fill="none"/>
            <circle cx="9" cy="9" r="3" fill="currentColor"/>
          </svg>
        </template>
      </div>
      <div class="bubble-wrapper">
        <div class="bubble">
          <div v-if="isUser" class="content">{{ message.content }}</div>
          <div v-else class="content markdown-body" v-html="renderedContent"></div>
        </div>
        <div class="meta">
          <span class="role">{{ isUser ? '你' : 'AI' }}</span>
          <span class="time">{{ timeStr }}</span>
        </div>
      </div>
    </template>
  </div>
</template>

<style scoped>
.message {
  display: flex;
  gap: 12px;
  animation: fadeIn 0.3s ease-out;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(8px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.message-user {
  flex-direction: row-reverse;
}

.message-system {
  justify-content: center;
}

.system-msg {
  font-size: 12px;
  color: var(--text-muted);
  background: var(--bg-tertiary);
  padding: 4px 14px;
  border-radius: 20px;
  border: 1px solid var(--border);
}

/* Avatar */
.avatar {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  margin-top: 4px;
}

.avatar-user {
  background: var(--user-bubble);
  color: var(--accent-light);
  border: 1px solid var(--border-light);
}

.avatar-ai {
  background: linear-gradient(135deg, var(--accent), var(--accent-light));
  color: white;
}

/* Bubble */
.bubble-wrapper {
  max-width: 70%;
  min-width: 60px;
}

.message-user .bubble-wrapper {
  align-items: flex-end;
}

.bubble {
  padding: 12px 16px;
  border-radius: var(--radius-lg);
  line-height: 1.65;
  font-size: 14.5px;
  word-break: break-word;
}

/* User messages: preserve line breaks */
.message-user .content {
  white-space: pre-wrap;
}

.message-user .bubble {
  background: var(--user-bubble);
  color: var(--text-primary);
  border: 1px solid var(--border-light);
  border-bottom-right-radius: var(--radius-sm);
}

.message-ai .bubble {
  background: var(--ai-bubble);
  color: var(--text-primary);
  border: 1px solid var(--border);
  border-bottom-left-radius: var(--radius-sm);
}

.meta {
  display: flex;
  gap: 8px;
  margin-top: 6px;
  font-size: 11px;
  color: var(--text-muted);
}

.message-user .meta {
  justify-content: flex-end;
}

.role {
  font-weight: 500;
}

@media (max-width: 640px) {
  .bubble-wrapper {
    max-width: 85%;
  }
  .avatar {
    width: 30px;
    height: 30px;
  }
  .bubble {
    padding: 10px 14px;
    font-size: 14px;
  }
}

/* Markdown styles */
.markdown-body {
  line-height: 1.7;
}

.markdown-body :first-child {
  margin-top: 0;
}

.markdown-body :last-child {
  margin-bottom: 0;
}

.markdown-body h1,
.markdown-body h2,
.markdown-body h3,
.markdown-body h4 {
  margin: 16px 0 10px;
  font-weight: 600;
  line-height: 1.35;
}

.markdown-body h1 {
  font-size: 18px;
  border-bottom: 1px solid var(--border);
  padding-bottom: 6px;
}

.markdown-body h2 {
  font-size: 16px;
}

.markdown-body h3 {
  font-size: 15px;
}

.markdown-body p {
  margin: 10px 0;
}

.markdown-body strong {
  font-weight: 600;
  color: var(--text-primary);
}

.markdown-body em {
  font-style: italic;
}

.markdown-body a {
  color: var(--accent-light);
  text-decoration: none;
}

.markdown-body a:hover {
  text-decoration: underline;
}

.markdown-body ul,
.markdown-body ol {
  margin: 10px 0;
  padding-left: 20px;
}

.markdown-body li {
  margin: 4px 0;
}

.markdown-body ul li {
  list-style-type: disc;
}

.markdown-body ol li {
  list-style-type: decimal;
}

.markdown-body code {
  background: rgba(108, 92, 231, 0.15);
  color: var(--accent-light);
  padding: 2px 6px;
  border-radius: 4px;
  font-family: 'JetBrains Mono', 'Fira Code', 'SF Mono', Consolas, monospace;
  font-size: 0.9em;
}

.markdown-body pre {
  background: var(--bg-tertiary);
  border: 1px solid var(--border);
  border-radius: var(--radius-md);
  padding: 14px 16px;
  margin: 12px 0;
  overflow-x: auto;
}

.markdown-body pre code {
  background: none;
  color: var(--text-primary);
  padding: 0;
  font-size: 13px;
  line-height: 1.6;
}

.markdown-body blockquote {
  border-left: 3px solid var(--accent);
  margin: 12px 0;
  padding: 4px 14px;
  color: var(--text-secondary);
  background: rgba(108, 92, 231, 0.06);
  border-radius: 0 var(--radius-sm) var(--radius-sm) 0;
}

.markdown-body hr {
  border: none;
  border-top: 1px solid var(--border);
  margin: 16px 0;
}

.markdown-body table {
  width: 100%;
  border-collapse: collapse;
  margin: 12px 0;
  font-size: 13px;
}

.markdown-body th,
.markdown-body td {
  border: 1px solid var(--border);
  padding: 8px 12px;
  text-align: left;
}

.markdown-body th {
  background: var(--bg-tertiary);
  font-weight: 600;
}

.markdown-body tr:nth-child(even) {
  background: rgba(255, 255, 255, 0.02);
}

.markdown-body img {
  max-width: 100%;
  border-radius: var(--radius-md);
}
</style>
