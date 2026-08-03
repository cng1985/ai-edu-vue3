<script setup>
import { ref, nextTick, watch, onMounted } from 'vue'
import { useChatStore } from '../stores/chat'
import { useAuthStore } from '../stores/auth'
import MarkdownRenderer from '../components/MarkdownRenderer.vue'

const chat = useChatStore()
const auth = useAuthStore()
const input = ref('')
const scrollArea = ref(null)

const suggestions = [
  '什么是 RAG？它解决什么问题？',
  '如何防止 Agent 陷入工具调用死循环？',
  'Function Calling 的工作原理是什么？',
  '我该从哪门课开始学？'
]

const modeLabel = () => {
  if (!chat.aiConfig) return '连接中…'
  return chat.aiConfig.enabled
    ? `大模型 · ${chat.aiConfig.model}`
    : '大模型未配置'
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
      <div>
        <h1>AI 学习助手</h1>
        <p>
          基于平台课程知识库的 RAG 智能答疑，支持多轮对话、流式输出与出处溯源。
          <span class="chat__mode" :class="{ 'chat__mode--live': chat.aiConfig?.enabled }">{{ modeLabel() }}</span>
        </p>
      </div>
      <button v-if="chat.messages.length" class="btn btn--danger" @click="chat.clear()">
        清空对话
      </button>
    </header>

    <div ref="scrollArea" class="chat__scroll">
      <div v-if="chat.messages.length === 0" class="chat__empty">
        <div class="chat__empty-icon">🤖</div>
        <h2>你好，我是你的 AI 学习助手</h2>
        <p>我熟悉平台的全部课程内容，可以帮你解答概念、定位知识点出处、推荐学习路径。</p>
        <div class="chat__suggestions">
          <button
            v-for="s in suggestions"
            :key="s"
            class="chat__suggestion"
            @click="send(s)"
          >
            {{ s }}
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
          :style="
            msg.role === 'user'
              ? {
                  background: auth.user?.avatarColor || '#1772f6',
                  color: '#fff',
                  borderColor: 'transparent'
                }
              : undefined
          "
        >
          {{ msg.role === 'user' ? auth.user?.avatar || '我' : 'AI' }}
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
              <span class="chat__sources-label">📎 知识来源：</span>
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
      <textarea
        v-model="input"
        rows="1"
        placeholder="输入你的问题，Enter 发送，Shift+Enter 换行…"
        @keydown="onKeydown"
      ></textarea>
      <button
        v-if="chat.generating"
        class="btn btn--danger"
        @click="chat.stop()"
      >
        ⏹ 停止
      </button>
      <button v-else class="btn btn--primary" :disabled="!input.trim()" @click="send()">
        发送 ↑
      </button>
    </footer>
  </div>
</template>

<style scoped>
.chat {
  display: flex;
  flex-direction: column;
  height: 100vh;
  max-width: 900px;
  margin: 0 auto;
  padding: 0 32px;
}

.chat__header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 16px;
  padding: 24px 4px 16px;
  border-bottom: 1px solid var(--border);
}

.chat__header h1 {
  margin: 0 0 3px;
  font-size: 21px;
}

.chat__header p {
  margin: 0;
  font-size: 13px;
  color: var(--text-2);
}

.chat__mode {
  display: inline-block;
  margin-left: 6px;
  padding: 2px 8px;
  border-radius: 999px;
  background: var(--surface-2);
  color: var(--text-3);
  font-size: 11px;
  font-weight: 600;
}

.chat__mode--live {
  background: #ecfdf5;
  color: #047857;
}

.chat__scroll {
  flex: 1;
  overflow-y: auto;
  padding: 22px 4px;
  display: flex;
  flex-direction: column;
  gap: 18px;
}

.chat__empty {
  margin: auto;
  text-align: center;
  max-width: 480px;
  padding: 24px 0;
}

.chat__empty-icon {
  font-size: 46px;
  margin-bottom: 8px;
}

.chat__empty h2 {
  margin: 0 0 8px;
  font-size: 19px;
}

.chat__empty p {
  margin: 0 0 20px;
  color: var(--text-2);
  font-size: 14px;
}

.chat__suggestions {
  display: flex;
  flex-direction: column;
  gap: 9px;
}

.chat__suggestion {
  padding: 11px 16px;
  border: 1px solid var(--border);
  background: var(--surface);
  border-radius: var(--radius-sm);
  font-size: 13.5px;
  color: var(--text-2);
  cursor: pointer;
  text-align: left;
  transition: all 0.15s ease;
}

.chat__suggestion:hover {
  border-color: var(--primary);
  color: var(--primary);
  background: var(--primary-soft);
}

.chat__row {
  display: flex;
  gap: 11px;
  align-items: flex-start;
}

.chat__row--user {
  flex-direction: row-reverse;
}

.chat__avatar {
  width: 36px;
  height: 36px;
  min-width: 36px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  background: var(--surface);
  border: 1px solid var(--border);
  font-size: 12px;
  font-weight: 700;
  color: #1772f6;
}

.chat__bubble {
  max-width: 78%;
  padding: 12px 17px;
  border-radius: 16px;
  font-size: 14.5px;
  line-height: 1.7;
}

.chat__bubble--user {
  background: var(--primary);
  color: #fff;
  border-top-right-radius: 5px;
  white-space: pre-wrap;
}

.chat__bubble--assistant {
  background: var(--surface);
  border: 1px solid var(--border);
  border-top-left-radius: 5px;
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
  background: var(--text-3);
  animation: bounce 1.2s infinite ease-in-out;
}

.chat__thinking span:nth-child(2) { animation-delay: 0.15s; }
.chat__thinking span:nth-child(3) { animation-delay: 0.3s; }

@keyframes bounce {
  0%, 70%, 100% { transform: translateY(0); opacity: 0.45; }
  35% { transform: translateY(-5px); opacity: 1; }
}

.chat__sources {
  margin-top: 12px;
  padding-top: 10px;
  border-top: 1px dashed var(--border);
  display: flex;
  flex-wrap: wrap;
  gap: 7px;
  align-items: center;
}

.chat__sources-label {
  font-size: 12px;
  color: var(--text-3);
}

.chat__source {
  font-size: 12px;
  padding: 3px 10px;
  border-radius: 999px;
  background: var(--primary-soft);
  color: var(--primary-strong);
  transition: all 0.15s ease;
}

.chat__source:hover {
  background: var(--primary);
  color: #fff;
}

.chat__composer {
  display: flex;
  gap: 11px;
  align-items: flex-end;
  padding: 14px 4px 22px;
  border-top: 1px solid var(--border);
}

.chat__composer textarea {
  flex: 1;
  border: 1px solid var(--border);
  border-radius: var(--radius-sm);
  padding: 11px 15px;
  font-family: inherit;
  font-size: 14.5px;
  line-height: 1.6;
  resize: none;
  outline: none;
  max-height: 140px;
  background: var(--surface);
  transition: border-color 0.15s ease;
}

.chat__composer textarea:focus {
  border-color: var(--primary);
}

@media (max-width: 720px) {
  .chat {
    padding: 0 14px;
    height: calc(100vh - 52px);
  }

  .chat__bubble {
    max-width: 88%;
  }
}
</style>
