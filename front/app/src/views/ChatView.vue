<script setup>
import { ref, nextTick, watch, onMounted } from 'vue'
import { useChatStore } from '../stores/chat'
import { useAuthStore } from '../stores/auth'
import MarkdownRenderer from '../components/MarkdownRenderer.vue'
import Icon from '../components/Icon.vue'

const chat = useChatStore()
const auth = useAuthStore()
const input = ref('')
const scrollArea = ref(null)

const suggestions = [
  { icon: 'layers', text: '什么是 RAG？它解决什么问题？' },
  { icon: 'refresh', text: '如何防止 Agent 陷入工具调用死循环？' },
  { icon: 'code', text: 'Function Calling 的工作原理是什么？' },
  { icon: 'compass', text: '我该从哪门课开始学？' }
]

const modeLabel = () => {
  if (!chat.aiConfig) return '连接中…'
  return chat.aiConfig.enabled
    ? `大模型 · ${chat.aiConfig.model}`
    : '本地知识库模式'
}

function send(text) {
  const question = (text ?? input.value).trim()
  if (!question) return
  chat.send(question)
  input.value = ''
}

function onKeydown(e) {
  if (e.key === 'Enter' && !e.shiftKey && !e.isComposing) {
    e.preventDefault()
    send()
  }
}

async function scrollToBottom() {
  await nextTick()
  if (scrollArea.value) {
    scrollArea.value.scrollTop = scrollArea.value.scrollHeight
  }
}

watch(
  () => chat.messages.map((m) => m.text.length).join(','),
  scrollToBottom
)
onMounted(() => {
  chat.loadConfig()
  scrollToBottom()
})
</script>

<template>
  <div class="chat">
    <header class="chat__header">
      <div class="chat__title">
        <span class="chat__logo"><Icon name="sparkles" :size="20" /></span>
        <div>
          <h1>AI 学习助手</h1>
          <p>
            基于课程知识库的 RAG 智能答疑 · 多轮对话 · 流式输出 · 出处溯源
          </p>
        </div>
      </div>
      <div class="chat__actions">
        <span class="chat__mode" :class="{ 'chat__mode--live': chat.aiConfig?.enabled }">
          <i></i>{{ modeLabel() }}
        </span>
        <button v-if="chat.messages.length" class="btn btn--ghost btn--sm" @click="chat.clear()">
          <Icon name="trash" :size="14" /> 清空对话
        </button>
      </div>
    </header>

    <div ref="scrollArea" class="chat__scroll">
      <div v-if="chat.messages.length === 0" class="chat__empty">
        <div class="chat__empty-icon"><Icon name="sparkles" :size="34" /></div>
        <h2>你好，我是你的 AI 学习助手</h2>
        <p>我熟悉平台的全部课程内容，可以帮你解答概念、定位知识点出处、推荐学习路径。</p>
        <div class="chat__suggestions">
          <button
            v-for="s in suggestions"
            :key="s.text"
            class="chat__suggestion"
            @click="send(s.text)"
          >
            <span class="chat__suggestion-icon"><Icon :name="s.icon" :size="16" /></span>
            <span>{{ s.text }}</span>
            <Icon name="arrowRight" :size="14" class="chat__suggestion-arrow" />
          </button>
        </div>
      </div>

      <div
        v-for="msg in chat.messages"
        :key="msg.id"
        class="chat__row"
        :class="`chat__row--${msg.role}`"
      >
        <div
          class="chat__avatar"
          :class="{ 'chat__avatar--user': msg.role === 'user' }"
          :style="msg.role === 'user' ? { background: auth.user?.avatarColor || '#6b5cff' } : undefined"
        >
          <template v-if="msg.role === 'user'">{{ auth.user?.avatar || '我' }}</template>
          <Icon v-else name="sparkles" :size="16" />
        </div>
        <div class="chat__bubble" :class="`chat__bubble--${msg.role}`">
          <template v-if="msg.role === 'assistant'">
            <div v-if="msg.streaming && !msg.text" class="chat__thinking">
              <span></span><span></span><span></span>
            </div>
            <MarkdownRenderer
              v-else
              :source="msg.text"
              :live="msg.streaming"
              :class="{ 'cursor-blink': msg.streaming }"
            />
            <div v-if="msg.sources.length" class="chat__sources">
              <span class="chat__sources-label"><Icon name="link" :size="12" /> 知识来源</span>
              <router-link
                v-for="src in msg.sources"
                :key="src.courseId + src.chapterId"
                :to="`/courses/${src.courseId}/${src.chapterId}`"
                class="chat__source"
              >
                {{ src.courseTitle }} · {{ src.chapterTitle }}
              </router-link>
            </div>
          </template>
          <template v-else>{{ msg.text }}</template>
        </div>
      </div>
    </div>

    <footer class="chat__composer">
      <div class="chat__composer-box">
        <textarea
          v-model="input"
          rows="1"
          placeholder="输入你的问题，Enter 发送，Shift+Enter 换行…"
          @keydown="onKeydown"
        ></textarea>
        <button
          v-if="chat.generating"
          class="chat__send chat__send--stop"
          title="停止生成"
          @click="chat.stop()"
        >
          <Icon name="stop" :size="16" />
        </button>
        <button v-else class="chat__send" :disabled="!input.trim()" title="发送" @click="send()">
          <Icon name="send" :size="16" />
        </button>
      </div>
      <p class="chat__disclaimer">AI 生成内容可能存在偏差，请结合课程原文与实践验证。</p>
    </footer>
  </div>
</template>

<style scoped>
.chat {
  display: flex;
  flex-direction: column;
  height: 100vh;
  max-width: 920px;
  margin: 0 auto;
  padding: 0 32px;
}

.chat__header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 16px;
  padding: 26px 4px 18px;
  border-bottom: 1px solid var(--border);
}

.chat__title {
  display: flex;
  align-items: center;
  gap: 14px;
}

.chat__logo {
  display: grid;
  place-items: center;
  width: 44px;
  height: 44px;
  border-radius: 14px;
  color: #fff;
  background: linear-gradient(135deg, var(--primary), #a071ff);
  box-shadow: var(--shadow-primary);
}

.chat__header h1 {
  margin: 0 0 2px;
  font-size: 20px;
}

.chat__header p {
  margin: 0;
  font-size: 13px;
  color: var(--text-2);
}

.chat__actions {
  display: flex;
  align-items: center;
  gap: 10px;
}

.chat__mode {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 5px 11px;
  border-radius: 999px;
  background: var(--surface-3);
  color: var(--text-2);
  font-size: 12px;
  font-weight: 600;
  white-space: nowrap;
}

.chat__mode i {
  width: 7px;
  height: 7px;
  border-radius: 50%;
  background: var(--text-3);
}

.chat__mode--live {
  background: var(--success-soft);
  color: var(--success-strong);
}

.chat__mode--live i {
  background: var(--success);
  box-shadow: 0 0 0 3px rgba(15, 185, 129, 0.2);
}

.chat__scroll {
  flex: 1;
  overflow-y: auto;
  padding: 24px 4px;
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.chat__empty {
  margin: auto;
  text-align: center;
  max-width: 520px;
  padding: 24px 0;
}

.chat__empty-icon {
  display: grid;
  place-items: center;
  width: 72px;
  height: 72px;
  margin: 0 auto 16px;
  border-radius: 24px;
  color: #fff;
  background: linear-gradient(135deg, var(--primary), #a071ff);
  box-shadow: var(--shadow-primary);
  animation: float 4s ease-in-out infinite;
}

.chat__empty h2 {
  margin: 0 0 8px;
  font-size: 22px;
}

.chat__empty p {
  margin: 0 0 24px;
  color: var(--text-2);
  font-size: 14.5px;
}

.chat__suggestions {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 10px;
}

.chat__suggestion {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 14px;
  border: 1px solid var(--border);
  background: var(--surface);
  border-radius: 14px;
  font: inherit;
  font-size: 13.5px;
  color: var(--text-2);
  cursor: pointer;
  text-align: left;
  transition: all var(--t-fast);
}

.chat__suggestion > span:nth-child(2) {
  flex: 1;
}

.chat__suggestion-icon {
  display: grid;
  place-items: center;
  width: 30px;
  height: 30px;
  flex: 0 0 30px;
  border-radius: 9px;
  background: var(--primary-soft);
  color: var(--primary-strong);
}

.chat__suggestion-arrow {
  opacity: 0;
  transform: translateX(-4px);
  transition: all var(--t-fast);
}

.chat__suggestion:hover {
  border-color: var(--primary);
  color: var(--text);
  box-shadow: var(--shadow-sm);
  transform: translateY(-1px);
}

.chat__suggestion:hover .chat__suggestion-arrow {
  opacity: 1;
  transform: translateX(0);
  color: var(--primary);
}

.chat__row {
  display: flex;
  gap: 12px;
  align-items: flex-start;
  animation: fade-up 0.3s var(--ease) both;
}

.chat__row--user {
  flex-direction: row-reverse;
}

.chat__avatar {
  width: 36px;
  height: 36px;
  min-width: 36px;
  display: grid;
  place-items: center;
  border-radius: 12px;
  background: linear-gradient(135deg, var(--primary), #a071ff);
  color: #fff;
  font-size: 12px;
  font-weight: 700;
  box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.3);
}

.chat__bubble {
  max-width: 78%;
  padding: 13px 18px;
  border-radius: 18px;
  font-size: 14.5px;
  line-height: 1.7;
}

.chat__bubble--user {
  background: linear-gradient(135deg, var(--primary), var(--primary-strong));
  color: #fff;
  border-top-right-radius: 6px;
  white-space: pre-wrap;
  box-shadow: 0 6px 16px var(--primary-glow);
}

.chat__bubble--assistant {
  background: var(--surface);
  border: 1px solid var(--border);
  border-top-left-radius: 6px;
  box-shadow: var(--shadow-xs);
}

.chat__bubble--assistant :deep(.markdown-body) {
  font-size: 14.5px;
}

.chat__bubble--assistant :deep(.markdown-body > :first-child) {
  margin-top: 0;
}

.chat__bubble--assistant :deep(.markdown-body > :last-child) {
  margin-bottom: 0;
}

.chat__thinking {
  display: flex;
  gap: 5px;
  padding: 6px 2px;
}

.chat__thinking span {
  width: 7px;
  height: 7px;
  border-radius: 50%;
  background: var(--primary);
  animation: bounce 1.2s infinite ease-in-out;
}

.chat__thinking span:nth-child(2) { animation-delay: 0.15s; }
.chat__thinking span:nth-child(3) { animation-delay: 0.3s; }

@keyframes bounce {
  0%, 70%, 100% { transform: translateY(0); opacity: 0.35; }
  35% { transform: translateY(-5px); opacity: 1; }
}

.chat__sources {
  margin-top: 14px;
  padding-top: 12px;
  border-top: 1px dashed var(--border);
  display: flex;
  flex-wrap: wrap;
  gap: 7px;
  align-items: center;
}

.chat__sources-label {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  color: var(--text-3);
  font-weight: 600;
}

.chat__source {
  font-size: 12px;
  padding: 4px 11px;
  border-radius: 999px;
  background: var(--primary-soft);
  color: var(--primary-strong);
  font-weight: 500;
  transition: all var(--t-fast);
}

.chat__source:hover {
  background: var(--primary);
  color: #fff;
}

.chat__composer {
  padding: 14px 4px 20px;
  border-top: 1px solid var(--border);
}

.chat__composer-box {
  display: flex;
  gap: 10px;
  align-items: flex-end;
  padding: 8px 8px 8px 16px;
  border: 1px solid var(--border-strong);
  border-radius: 18px;
  background: var(--surface);
  box-shadow: var(--shadow-sm);
  transition: border-color var(--t-fast), box-shadow var(--t-fast);
}

.chat__composer-box:focus-within {
  border-color: var(--primary);
  box-shadow: 0 0 0 4px var(--primary-soft-2);
}

.chat__composer textarea {
  flex: 1;
  border: none;
  padding: 8px 0;
  font-family: inherit;
  font-size: 14.5px;
  line-height: 1.6;
  resize: none;
  outline: none;
  max-height: 160px;
  background: transparent;
  color: var(--text);
}

.chat__composer textarea::placeholder {
  color: var(--text-3);
}

.chat__send {
  display: grid;
  place-items: center;
  width: 40px;
  height: 40px;
  flex: 0 0 40px;
  border: none;
  border-radius: 12px;
  background: linear-gradient(135deg, var(--primary), var(--primary-strong));
  color: #fff;
  cursor: pointer;
  box-shadow: 0 6px 14px var(--primary-glow);
  transition: transform var(--t-fast), box-shadow var(--t-fast), opacity var(--t-fast);
}

.chat__send:hover:not(:disabled) {
  transform: translateY(-1px);
}

.chat__send:disabled {
  opacity: 0.4;
  box-shadow: none;
  cursor: not-allowed;
}

.chat__send--stop {
  background: var(--danger);
  box-shadow: 0 6px 14px rgba(239, 77, 99, 0.3);
}

.chat__disclaimer {
  margin: 10px 0 0;
  text-align: center;
  color: var(--text-3);
  font-size: 11.5px;
}

@media (max-width: 720px) {
  .chat {
    padding: 0 14px;
    height: calc(100vh - 56px);
  }

  .chat__header {
    flex-direction: column;
    align-items: flex-start;
  }

  .chat__bubble {
    max-width: 88%;
  }

  .chat__suggestions {
    grid-template-columns: 1fr;
  }
}
</style>
