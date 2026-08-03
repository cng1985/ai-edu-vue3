<template>
  <div class="cgpt" :class="{ 'cgpt--sidebar-collapsed': sidebarCollapsed }">
    <!-- 左侧边栏 -->
    <aside class="cgpt-sidebar" :class="{ 'cgpt-sidebar--hidden': sidebarCollapsed }">
      <div class="cgpt-sidebar__inner">
        <button type="button" class="cgpt-btn-new" @click="onNewChat">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M12 5v14M5 12h14" stroke-linecap="round" />
          </svg>
          <span>新聊天</span>
        </button>

        <div class="cgpt-sidebar__list">
          <button
            v-for="s in sessions"
            :key="s.id"
            type="button"
            class="cgpt-history-item"
            :class="{ 'cgpt-history-item--active': s.id === activeId }"
            @click="switchSession(s.id)"
          >
            <span class="cgpt-history-item__text">{{ s.title }}</span>
            <span
              class="cgpt-history-item__delete"
              title="删除"
              role="button"
              @click.stop="onDeleteSession(s.id)"
            >
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M18 6L6 18M6 6l12 12" stroke-linecap="round" />
              </svg>
            </span>
          </button>
        </div>

        <div class="cgpt-sidebar__footer">
          <button type="button" class="cgpt-sidebar-user" @click="router.push('/dashboard')">
            <span class="cgpt-sidebar-user__avatar" :style="{ background: auth.user?.avatarColor || '#10a37f' }">
              {{ auth.user?.avatar || '管' }}
            </span>
            <span class="cgpt-sidebar-user__name">{{ auth.user?.nickname || '管理员' }}</span>
          </button>
        </div>
      </div>
    </aside>

    <!-- 主区域 -->
    <main class="cgpt-main">
      <header class="cgpt-header">
        <button type="button" class="cgpt-icon-btn" title="切换侧边栏" @click="sidebarCollapsed = !sidebarCollapsed">
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <rect x="3" y="4" width="18" height="16" rx="2" />
            <path d="M9 4v16" />
          </svg>
        </button>

        <div class="cgpt-model-picker" ref="modelPickerRef">
          <button type="button" class="cgpt-model-btn" @click="modelMenuOpen = !modelMenuOpen">
            <span>{{ currentModelLabel }}</span>
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M6 9l6 6 6-6" stroke-linecap="round" stroke-linejoin="round" />
            </svg>
          </button>
          <div v-if="modelMenuOpen" class="cgpt-model-menu">
            <button
              v-for="vm in virtualModels"
              :key="vm.code"
              type="button"
              class="cgpt-model-option"
              :class="{ 'cgpt-model-option--active': virtualModel === vm.code }"
              @click="selectModel(vm.code)"
            >
              {{ vm.name }}
            </button>
            <div class="cgpt-model-menu__divider" />
            <button
              type="button"
              class="cgpt-model-option"
              :class="{ 'cgpt-model-option--active': mode === 'rag' }"
              @click="setMode('rag')"
            >
              知识库增强模式
            </button>
            <button
              type="button"
              class="cgpt-model-option"
              :class="{ 'cgpt-model-option--active': mode === 'chat' }"
              @click="setMode('chat')"
            >
              纯对话模式
            </button>
          </div>
        </div>

        <div class="cgpt-header__spacer" />
      </header>

      <div ref="scrollArea" class="cgpt-body">
        <!-- 空状态 -->
        <div v-if="!messages.length" class="cgpt-welcome">
          <h1 class="cgpt-welcome__title">有什么可以帮忙的？</h1>
          <div class="cgpt-welcome__grid">
            <button
              v-for="(item, i) in welcomeCards"
              :key="i"
              type="button"
              class="cgpt-welcome-card"
              @click="send(item.prompt)"
            >
              <span class="cgpt-welcome-card__title">{{ item.title }}</span>
              <span class="cgpt-welcome-card__desc">{{ item.desc }}</span>
            </button>
          </div>
        </div>

        <!-- 消息列表 -->
        <template v-else>
          <article
            v-for="(msg, index) in messages"
            :key="msg.id"
            class="cgpt-msg"
            :class="`cgpt-msg--${msg.role}`"
          >
            <div class="cgpt-msg__inner">
              <div class="cgpt-msg__avatar">
                <template v-if="msg.role === 'assistant'">
                  <svg class="cgpt-logo" viewBox="0 0 41 41" fill="none">
                    <path fill="#10a37f" d="M37.5 18.5c0-1.2-.3-2.4-.9-3.5L33 9.2c-.6-1.1-1.5-2-2.6-2.6L26.5 3.4C25.4 2.8 24.2 2.5 23 2.5h-4c-1.2 0-2.4.3-3.5.9L9.2 6.9C8.1 7.5 7.2 8.4 6.6 9.5L3.4 15.5c-.6 1.1-.9 2.3-.9 3.5v4c0 1.2.3 2.4.9 3.5l3.2 6c.6 1.1 1.5 2 2.6 2.6l6 3.2c1.1.6 2.3.9 3.5.9h4c1.2 0 2.4-.3 3.5-.9l6-3.2c1.1-.6 2-1.5 2.6-2.6l3.2-6c.6-1.1.9-2.3.9-3.5v-4z" />
                    <path fill="#fff" d="M20.5 11v9.5l6 3.5" stroke="#fff" stroke-width="1.5" />
                  </svg>
                </template>
                <template v-else>
                  <span class="cgpt-msg__user-icon" :style="{ background: auth.user?.avatarColor || '#5436DA' }">
                    {{ auth.user?.avatar || '我' }}
                  </span>
                </template>
              </div>
              <div class="cgpt-msg__content">
                <div v-if="msg.streaming && !msg.text" class="cgpt-dots">
                  <span /><span /><span />
                </div>
                <ChatGptMarkdown v-else :source="msg.text" :live="msg.streaming" />
                <div v-if="msg.role === 'assistant' && !msg.streaming && msg.text" class="cgpt-msg__actions">
                  <button type="button" class="cgpt-msg-action" title="复制" @click="copyText(msg.text)">
                    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <rect x="9" y="9" width="13" height="13" rx="2" />
                      <path d="M5 15H4a2 2 0 01-2-2V4a2 2 0 012-2h9a2 2 0 012 2v1" />
                    </svg>
                  </button>
                  <button
                    v-if="index === messages.length - 1 && !generating"
                    type="button"
                    class="cgpt-msg-action"
                    title="重新生成"
                    @click="regenerate"
                  >
                    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <path d="M1 4v6h6M23 20v-6h-6" stroke-linecap="round" stroke-linejoin="round" />
                      <path d="M20.49 9A9 9 0 005.64 5.64L1 10m22 4l-4.64 4.36A9 9 0 013.51 15" stroke-linecap="round" stroke-linejoin="round" />
                    </svg>
                  </button>
                </div>
              </div>
            </div>
          </article>
        </template>
      </div>

      <!-- 底部输入 -->
      <footer class="cgpt-footer">
        <div class="cgpt-input-box">
          <div class="cgpt-input-shell" :class="{ 'cgpt-input-shell--focus': inputFocused }">
            <textarea
              ref="inputRef"
              v-model="input"
              class="cgpt-input"
              rows="1"
              placeholder="询问任何问题"
              :disabled="generating"
              @focus="inputFocused = true"
              @blur="inputFocused = false"
              @input="resizeInput"
              @keydown="onKeydown"
            />
            <div class="cgpt-input-actions">
              <button
                v-if="generating"
                type="button"
                class="cgpt-send cgpt-send--stop"
                title="停止生成"
                @click="stop"
              >
                <svg width="14" height="14" viewBox="0 0 24 24" fill="currentColor">
                  <rect x="6" y="6" width="12" height="12" rx="1" />
                </svg>
              </button>
              <button
                v-else
                type="button"
                class="cgpt-send"
                :class="{ 'cgpt-send--active': canSend }"
                :disabled="!canSend"
                title="发送"
                @click="send()"
              >
                <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
                  <path d="M12 19V5M5 12l7-7 7 7" stroke-linecap="round" stroke-linejoin="round" />
                </svg>
              </button>
            </div>
          </div>
          <p class="cgpt-disclaimer">ChatGPT 可能会犯错。请核查重要信息。</p>
        </div>
      </footer>
    </main>
  </div>
</template>

<script setup>
import { computed, nextTick, onMounted, onUnmounted, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessageBox } from 'element-plus'
import { useAuthStore } from '../stores/auth.js'
import { useChatGptSessions } from '../composables/useChatGptSessions.js'
import { aiApi } from '../api/ai.js'
import { aiModelsApi } from '../api/aiModels.js'
import ChatGptMarkdown from '../components/chatgpt/ChatGptMarkdown.vue'

const router = useRouter()
const auth = useAuthStore()
const {
  sessions,
  activeId,
  activeSession,
  createSession,
  selectSession,
  deleteSession,
  updateSession,
  setMessages
} = useChatGptSessions()

const sidebarCollapsed = ref(false)
const input = ref('')
const inputRef = ref(null)
const scrollArea = ref(null)
const inputFocused = ref(false)
const generating = ref(false)
const virtualModels = ref([])
const virtualModel = ref('chat-default')
const mode = ref('chat')
const modelMenuOpen = ref(false)
const modelPickerRef = ref(null)
let cancelStream = null

const messages = computed({
  get: () => activeSession.value?.messages || [],
  set: (val) => {
    if (activeSession.value) setMessages(activeSession.value.id, val)
  }
})

const canSend = computed(() => input.value.trim().length > 0 && !generating.value)

const currentModelLabel = computed(() => {
  const vm = virtualModels.value.find((v) => v.code === virtualModel.value)
  return vm?.name || virtualModel.value || 'ChatGPT'
})

const welcomeCards = [
  { title: '解释概念', desc: '用简单的话解释 RAG', prompt: '用简单的话解释什么是 RAG，以及它解决了什么问题' },
  { title: '写代码', desc: 'Python 快速排序', prompt: '用 Python 实现快速排序，并解释时间复杂度' },
  { title: '创意写作', desc: '产品介绍文案', prompt: '帮我写一段 100 字左右的 AI 学习平台产品介绍' },
  { title: '学习建议', desc: '如何入门后端', prompt: '零基础如何系统学习 Java 后端开发？给出学习路径' }
]

function persistSessionMeta() {
  if (!activeSession.value) return
  updateSession(activeSession.value.id, { virtualModel: virtualModel.value, mode: mode.value })
}

watch(activeId, () => {
  if (activeSession.value) {
    virtualModel.value = activeSession.value.virtualModel || virtualModel.value
    mode.value = activeSession.value.mode || 'chat'
  }
  nextTick(() => scrollToBottom())
})

function resizeInput() {
  const el = inputRef.value
  if (!el) return
  el.style.height = 'auto'
  el.style.height = `${Math.min(el.scrollHeight, 200)}px`
}

function scrollToBottom() {
  nextTick(() => {
    if (scrollArea.value) scrollArea.value.scrollTop = scrollArea.value.scrollHeight
  })
}

watch(() => messages.value.map((m) => m.text?.length).join(','), scrollToBottom)

async function loadMeta() {
  try {
    const cfg = await aiApi.config()
    if (cfg?.defaultVirtualModel) virtualModel.value = cfg.defaultVirtualModel
    virtualModels.value = await aiModelsApi.listVirtualModelOptions()
    if (!virtualModels.value.length) {
      virtualModels.value = [
        { code: 'chat-default', name: 'GPT-4o mini' },
        { code: 'chat-smart', name: 'GPT-4o' }
      ]
    }
  } catch {
    virtualModels.value = [
      { code: 'chat-default', name: 'GPT-4o mini' },
      { code: 'chat-smart', name: 'GPT-4o' }
    ]
  }
  if (activeSession.value && !activeSession.value.virtualModel) {
    updateSession(activeSession.value.id, { virtualModel: virtualModel.value, mode: mode.value })
  } else if (activeSession.value?.virtualModel) {
    virtualModel.value = activeSession.value.virtualModel
    mode.value = activeSession.value.mode || 'chat'
  }
}

function selectModel(code) {
  virtualModel.value = code
  modelMenuOpen.value = false
  persistSessionMeta()
}

function setMode(m) {
  mode.value = m
  modelMenuOpen.value = false
  persistSessionMeta()
}

function onNewChat() {
  stop()
  createSession()
  virtualModel.value = virtualModels.value[0]?.code || 'chat-default'
  mode.value = 'chat'
  persistSessionMeta()
}

function switchSession(id) {
  if (id === activeId.value) return
  stop()
  selectSession(id)
}

async function onDeleteSession(id) {
  try {
    await ElMessageBox.confirm('删除此对话？', '确认', { type: 'warning', confirmButtonText: '删除', cancelButtonText: '取消' })
    stop()
    deleteSession(id)
  } catch { /* cancelled */ }
}

function buildHistory(excludeLast = 1) {
  return messages.value
    .filter((m) => !m.streaming && m.text)
    .slice(0, -excludeLast)
    .slice(-16)
    .map((m) => ({ role: m.role, content: m.text }))
}

function runStream(question, history) {
  const reply = {
    id: `${Date.now()}:assistant`,
    role: 'assistant',
    text: '',
    streaming: true
  }
  const list = [...messages.value, reply]
  messages.value = list
  generating.value = true

  cancelStream = aiApi.chatStream(
    question,
    history,
    { virtualModel: virtualModel.value, mode: mode.value },
    {
      onToken: (chunk) => { reply.text += chunk },
      onDone: (payload) => {
        reply.text = payload.text || reply.text
        reply.streaming = false
        generating.value = false
        cancelStream = null
        messages.value = [...messages.value]
      },
      onError: (err) => {
        reply.text = `出错了：${err.message}`
        reply.streaming = false
        generating.value = false
        cancelStream = null
        messages.value = [...messages.value]
      }
    }
  )
}

function send(text) {
  const question = (text ?? input.value).trim()
  if (!question || generating.value) return

  const userMsg = {
    id: `${Date.now()}:user`,
    role: 'user',
    text: question,
    streaming: false
  }
  messages.value = [...messages.value, userMsg]
  input.value = ''
  resizeInput()
  persistSessionMeta()

  const history = buildHistory(0)
  runStream(question, history)
}

function regenerate() {
  const list = messages.value.filter((m) => !m.streaming)
  if (list.length < 2) return
  const lastUser = [...list].reverse().find((m) => m.role === 'user')
  if (!lastUser) return
  messages.value = list.slice(0, -1)
  const history = buildHistory(0)
  runStream(lastUser.text, history)
}

function stop() {
  if (cancelStream) {
    cancelStream()
    cancelStream = null
  }
  const list = messages.value
  const last = list[list.length - 1]
  if (last?.streaming) {
    last.streaming = false
    last.text += '\n\n*（已停止生成）*'
    messages.value = [...list]
  }
  generating.value = false
}

function onKeydown(e) {
  if (e.key === 'Enter' && !e.shiftKey && !e.isComposing) {
    e.preventDefault()
    if (canSend.value) send()
  }
}

async function copyText(text) {
  try {
    await navigator.clipboard.writeText(text)
  } catch { /* ignore */ }
}

function onDocClick(e) {
  if (modelMenuOpen.value && modelPickerRef.value && !modelPickerRef.value.contains(e.target)) {
    modelMenuOpen.value = false
  }
}

onMounted(() => {
  loadMeta()
  resizeInput()
  document.addEventListener('click', onDocClick)
})

onUnmounted(() => {
  stop()
  document.removeEventListener('click', onDocClick)
})
</script>

<style scoped>
/* ChatGPT 风格 — 深色侧栏 + 白色主区 */
.cgpt {
  display: flex;
  height: 100vh;
  width: 100%;
  background: #fff;
  color: #0d0d0d;
  font-family: ui-sans-serif, -apple-system, BlinkMacSystemFont, 'Segoe UI', Helvetica, Arial, sans-serif;
  overflow: hidden;
}

/* 侧栏 */
.cgpt-sidebar {
  width: 260px;
  flex-shrink: 0;
  background: #171717;
  color: #ececec;
  transition: width 0.2s ease, margin 0.2s ease;
}

.cgpt-sidebar--hidden {
  width: 0;
  overflow: hidden;
}

.cgpt-sidebar__inner {
  display: flex;
  flex-direction: column;
  height: 100%;
  padding: 8px;
  width: 260px;
}

.cgpt-btn-new {
  display: flex;
  align-items: center;
  gap: 10px;
  width: 100%;
  padding: 10px 12px;
  margin-bottom: 8px;
  border: 1px solid rgba(255, 255, 255, 0.15);
  border-radius: 12px;
  background: transparent;
  color: #ececec;
  font-size: 14px;
  cursor: pointer;
  transition: background 0.15s;
}

.cgpt-btn-new:hover {
  background: rgba(255, 255, 255, 0.08);
}

.cgpt-sidebar__list {
  flex: 1;
  overflow-y: auto;
  margin: 0 -4px;
  padding: 4px 0;
}

.cgpt-history-item {
  display: flex;
  align-items: center;
  gap: 4px;
  width: 100%;
  padding: 8px 10px;
  border: none;
  border-radius: 10px;
  background: transparent;
  color: #ececec;
  font-size: 14px;
  text-align: left;
  cursor: pointer;
  transition: background 0.15s;
}

.cgpt-history-item:hover,
.cgpt-history-item--active {
  background: rgba(255, 255, 255, 0.1);
}

.cgpt-history-item__text {
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.cgpt-history-item__delete {
  flex-shrink: 0;
  display: none;
  align-items: center;
  justify-content: center;
  width: 24px;
  height: 24px;
  border-radius: 6px;
  color: #9ca3af;
  cursor: pointer;
}

.cgpt-history-item:hover .cgpt-history-item__delete {
  display: flex;
}

.cgpt-history-item__delete:hover {
  background: rgba(255, 255, 255, 0.15);
  color: #fff;
}

.cgpt-sidebar__footer {
  padding-top: 8px;
  border-top: 1px solid rgba(255, 255, 255, 0.1);
}

.cgpt-sidebar-user {
  display: flex;
  align-items: center;
  gap: 10px;
  width: 100%;
  padding: 8px 10px;
  border: none;
  border-radius: 10px;
  background: transparent;
  color: #ececec;
  cursor: pointer;
  font-size: 14px;
}

.cgpt-sidebar-user:hover {
  background: rgba(255, 255, 255, 0.08);
}

.cgpt-sidebar-user__avatar {
  width: 28px;
  height: 28px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  color: #fff;
  flex-shrink: 0;
}

.cgpt-sidebar-user__name {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

/* 主区 */
.cgpt-main {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-width: 0;
  background: #fff;
}

.cgpt-header {
  display: flex;
  align-items: center;
  gap: 8px;
  height: 52px;
  padding: 0 12px;
  flex-shrink: 0;
}

.cgpt-icon-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  border: none;
  border-radius: 8px;
  background: transparent;
  color: #0d0d0d;
  cursor: pointer;
}

.cgpt-icon-btn:hover {
  background: #f4f4f4;
}

.cgpt-header__spacer {
  flex: 1;
}

.cgpt-model-picker {
  position: relative;
}

.cgpt-model-btn {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 6px 10px;
  border: none;
  border-radius: 10px;
  background: transparent;
  font-size: 18px;
  font-weight: 600;
  color: #0d0d0d;
  cursor: pointer;
}

.cgpt-model-btn:hover {
  background: #f4f4f4;
}

.cgpt-model-menu {
  position: absolute;
  top: calc(100% + 4px);
  left: 0;
  min-width: 200px;
  padding: 6px;
  background: #fff;
  border: 1px solid #e5e5e5;
  border-radius: 12px;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);
  z-index: 100;
}

.cgpt-model-option {
  display: block;
  width: 100%;
  padding: 10px 12px;
  border: none;
  border-radius: 8px;
  background: transparent;
  font-size: 14px;
  text-align: left;
  cursor: pointer;
  color: #0d0d0d;
}

.cgpt-model-option:hover {
  background: #f4f4f4;
}

.cgpt-model-option--active {
  background: #f0fdf4;
  color: #10a37f;
  font-weight: 500;
}

.cgpt-model-menu__divider {
  height: 1px;
  margin: 4px 8px;
  background: #e5e5e5;
}

/* 消息区 */
.cgpt-body {
  flex: 1;
  overflow-y: auto;
  scroll-behavior: smooth;
}

.cgpt-welcome {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 100%;
  padding: 40px 24px 120px;
}

.cgpt-welcome__title {
  font-size: 28px;
  font-weight: 600;
  margin: 0 0 32px;
  color: #0d0d0d;
}

.cgpt-welcome__grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px;
  max-width: 720px;
  width: 100%;
}

.cgpt-welcome-card {
  display: flex;
  flex-direction: column;
  gap: 4px;
  padding: 16px;
  border: 1px solid #e5e5e5;
  border-radius: 16px;
  background: #fff;
  text-align: left;
  cursor: pointer;
  transition: background 0.15s, border-color 0.15s;
}

.cgpt-welcome-card:hover {
  background: #f9f9f9;
  border-color: #d4d4d4;
}

.cgpt-welcome-card__title {
  font-size: 14px;
  font-weight: 600;
  color: #0d0d0d;
}

.cgpt-welcome-card__desc {
  font-size: 13px;
  color: #6b7280;
}

.cgpt-msg {
  border-bottom: 1px solid transparent;
}

.cgpt-msg--user {
  background: #f7f7f8;
}

.cgpt-msg__inner {
  display: flex;
  gap: 16px;
  max-width: 768px;
  margin: 0 auto;
  padding: 24px 24px;
}

.cgpt-msg__avatar {
  flex-shrink: 0;
  width: 32px;
  height: 32px;
}

.cgpt-logo {
  width: 32px;
  height: 32px;
}

.cgpt-msg__user-icon {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 13px;
  color: #fff;
}

.cgpt-msg__content {
  flex: 1;
  min-width: 0;
  padding-top: 2px;
}

.cgpt-msg__actions {
  display: flex;
  gap: 4px;
  margin-top: 8px;
}

.cgpt-msg-action {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border: none;
  border-radius: 8px;
  background: transparent;
  color: #6b7280;
  cursor: pointer;
}

.cgpt-msg-action:hover {
  background: #f4f4f4;
  color: #0d0d0d;
}

.cgpt-dots span {
  display: inline-block;
  width: 8px;
  height: 8px;
  margin-right: 4px;
  background: #9ca3af;
  border-radius: 50%;
  animation: cgpt-dot 1.2s infinite;
}

.cgpt-dots span:nth-child(2) { animation-delay: 0.15s; }
.cgpt-dots span:nth-child(3) { animation-delay: 0.3s; }

@keyframes cgpt-dot {
  0%, 80%, 100% { opacity: 0.3; transform: scale(0.85); }
  40% { opacity: 1; transform: scale(1); }
}

/* 底部输入 */
.cgpt-footer {
  flex-shrink: 0;
  padding: 12px 16px 20px;
  background: #fff;
}

.cgpt-input-box {
  max-width: 768px;
  margin: 0 auto;
}

.cgpt-input-shell {
  display: flex;
  align-items: flex-end;
  gap: 8px;
  padding: 10px 12px;
  border: 1px solid #d4d4d4;
  border-radius: 28px;
  background: #fff;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.06);
  transition: border-color 0.15s, box-shadow 0.15s;
}

.cgpt-input-shell--focus {
  border-color: #b4b4b4;
  box-shadow: 0 2px 16px rgba(0, 0, 0, 0.1);
}

.cgpt-input {
  flex: 1;
  border: none;
  outline: none;
  resize: none;
  font-size: 16px;
  line-height: 1.5;
  max-height: 200px;
  padding: 4px 0;
  font-family: inherit;
  background: transparent;
  color: #0d0d0d;
}

.cgpt-input::placeholder {
  color: #9ca3af;
}

.cgpt-input-actions {
  flex-shrink: 0;
}

.cgpt-send {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border: none;
  border-radius: 50%;
  background: #e5e5e5;
  color: #a3a3a3;
  cursor: not-allowed;
  transition: background 0.15s, color 0.15s;
}

.cgpt-send--active {
  background: #0d0d0d;
  color: #fff;
  cursor: pointer;
}

.cgpt-send--active:hover {
  background: #333;
}

.cgpt-send--stop {
  background: #0d0d0d;
  color: #fff;
  cursor: pointer;
}

.cgpt-disclaimer {
  margin: 10px 0 0;
  text-align: center;
  font-size: 12px;
  color: #9ca3af;
}

@media (max-width: 768px) {
  .cgpt-welcome__grid {
    grid-template-columns: 1fr;
  }

  .cgpt-sidebar--hidden + .cgpt-main {
    width: 100%;
  }
}
</style>
