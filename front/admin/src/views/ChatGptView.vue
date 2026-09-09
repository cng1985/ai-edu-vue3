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
            <span class="cgpt-sidebar-user__avatar" :style="{ background: auth.user?.avatarColor || 'var(--primary)' }">
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
                  <span class="cgpt-logo">AI</span>
                </template>
                <template v-else>
                  <span class="cgpt-msg__user-icon" :style="{ background: auth.user?.avatarColor || 'var(--primary)' }">
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
          <p class="cgpt-disclaimer">AI 可能会犯错，请核查重要信息。</p>
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
/* 对话工作台 — 深色侧栏 + 浅色主区，沿用全局设计系统 */
.cgpt {
  display: flex;
  height: 100vh;
  width: 100%;
  background: var(--bg);
  color: var(--text);
  overflow: hidden;
}

/* 侧栏 */
.cgpt-sidebar {
  width: 272px;
  flex-shrink: 0;
  background: linear-gradient(180deg, var(--ink-2) 0%, var(--ink) 100%);
  color: var(--ink-text);
  transition: width 0.22s var(--ease);
  position: relative;
}

.cgpt-sidebar::before {
  content: '';
  position: absolute;
  inset: 0;
  pointer-events: none;
  background: radial-gradient(320px 240px at 0% 0%, rgba(107, 92, 255, 0.35), transparent 70%);
}

.cgpt-sidebar--hidden {
  width: 0;
  overflow: hidden;
}

.cgpt-sidebar__inner {
  position: relative;
  display: flex;
  flex-direction: column;
  height: 100%;
  padding: 12px;
  width: 272px;
}

.cgpt-btn-new {
  display: flex;
  align-items: center;
  gap: 10px;
  width: 100%;
  padding: 11px 14px;
  margin-bottom: 10px;
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 13px;
  background: linear-gradient(135deg, rgba(107, 92, 255, 0.95) 0%, rgba(134, 118, 255, 0.85) 100%);
  color: #fff;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  box-shadow: 0 10px 24px rgba(107, 92, 255, 0.35);
  transition: transform var(--t-fast), box-shadow var(--t-fast);
}

.cgpt-btn-new:hover {
  transform: translateY(-1px);
  box-shadow: 0 14px 28px rgba(107, 92, 255, 0.45);
}

.cgpt-sidebar__list {
  flex: 1;
  overflow-y: auto;
  margin: 0 -4px;
  padding: 4px;
}

.cgpt-history-item {
  display: flex;
  align-items: center;
  gap: 4px;
  width: 100%;
  padding: 9px 12px;
  margin-bottom: 2px;
  border: none;
  border-radius: 11px;
  background: transparent;
  color: var(--ink-text);
  font-size: 13.5px;
  text-align: left;
  cursor: pointer;
  transition: background var(--t-fast), color var(--t-fast);
}

.cgpt-history-item:hover {
  background: rgba(255, 255, 255, 0.07);
  color: #fff;
}

.cgpt-history-item--active {
  background: rgba(255, 255, 255, 0.12);
  color: #fff;
  font-weight: 600;
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
  border-radius: 7px;
  color: var(--ink-text-2);
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
  padding-top: 10px;
  border-top: 1px solid rgba(255, 255, 255, 0.08);
}

.cgpt-sidebar-user {
  display: flex;
  align-items: center;
  gap: 10px;
  width: 100%;
  padding: 8px 10px;
  border: none;
  border-radius: 12px;
  background: transparent;
  color: var(--ink-text);
  cursor: pointer;
  font-size: 14px;
  transition: background var(--t-fast);
}

.cgpt-sidebar-user:hover {
  background: rgba(255, 255, 255, 0.07);
  color: #fff;
}

.cgpt-sidebar-user__avatar {
  width: 30px;
  height: 30px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: 700;
  color: #fff;
  flex-shrink: 0;
}

.cgpt-sidebar-user__name {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  font-weight: 600;
}

/* 主区 */
.cgpt-main {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-width: 0;
  background:
    radial-gradient(800px 400px at 85% -10%, rgba(107, 92, 255, 0.08), transparent 60%),
    var(--bg);
}

.cgpt-header {
  display: flex;
  align-items: center;
  gap: 10px;
  height: 60px;
  padding: 0 16px;
  flex-shrink: 0;
  background: rgba(255, 255, 255, 0.7);
  backdrop-filter: blur(16px);
  border-bottom: 1px solid var(--border);
}

.cgpt-icon-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 38px;
  height: 38px;
  border: 1px solid transparent;
  border-radius: 11px;
  background: transparent;
  color: var(--text-2);
  cursor: pointer;
  transition: all var(--t-fast);
}

.cgpt-icon-btn:hover {
  background: var(--surface-2);
  border-color: var(--border);
  color: var(--text);
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
  gap: 6px;
  padding: 8px 12px;
  border: 1px solid transparent;
  border-radius: 11px;
  background: transparent;
  font-family: 'Sora', 'Noto Sans SC', sans-serif;
  font-size: 15px;
  font-weight: 700;
  color: var(--text);
  cursor: pointer;
  transition: all var(--t-fast);
}

.cgpt-model-btn:hover {
  background: var(--surface-2);
  border-color: var(--border);
}

.cgpt-model-menu {
  position: absolute;
  top: calc(100% + 6px);
  left: 0;
  min-width: 240px;
  padding: 6px;
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: 14px;
  box-shadow: var(--shadow-lg);
  z-index: 100;
}

.cgpt-model-option {
  display: block;
  width: 100%;
  padding: 10px 12px;
  border: none;
  border-radius: 9px;
  background: transparent;
  font-size: 14px;
  font-weight: 500;
  text-align: left;
  cursor: pointer;
  color: var(--text-2);
  transition: all var(--t-fast);
}

.cgpt-model-option:hover {
  background: var(--surface-2);
  color: var(--text);
}

.cgpt-model-option--active {
  background: var(--primary-soft);
  color: var(--primary-deep);
  font-weight: 700;
}

.cgpt-model-menu__divider {
  height: 1px;
  margin: 6px 8px;
  background: var(--border);
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
  font-size: 30px;
  font-weight: 700;
  margin: 0 0 32px;
  color: var(--text);
  letter-spacing: -0.02em;
}

.cgpt-welcome__grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px;
  max-width: 760px;
  width: 100%;
}

.cgpt-welcome-card {
  display: flex;
  flex-direction: column;
  gap: 4px;
  padding: 16px 18px;
  border: 1px solid var(--border);
  border-radius: 16px;
  background: var(--surface);
  text-align: left;
  cursor: pointer;
  box-shadow: var(--shadow-sm);
  transition: all var(--t-fast);
}

.cgpt-welcome-card:hover {
  border-color: var(--primary);
  transform: translateY(-2px);
  box-shadow: var(--shadow);
}

.cgpt-welcome-card__title {
  font-size: 14px;
  font-weight: 700;
  color: var(--text);
}

.cgpt-welcome-card__desc {
  font-size: 13px;
  color: var(--text-3);
}

.cgpt-msg {
  border-bottom: 1px solid transparent;
}

.cgpt-msg__inner {
  display: flex;
  gap: 16px;
  max-width: 820px;
  margin: 0 auto;
  padding: 22px 24px;
}

.cgpt-msg__avatar {
  flex-shrink: 0;
  width: 34px;
  height: 34px;
}

.cgpt-logo {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 34px;
  height: 34px;
  border-radius: 11px;
  background: linear-gradient(135deg, var(--primary) 0%, #a396ff 100%);
  color: #fff;
  font-family: 'Sora', sans-serif;
  font-weight: 800;
  font-size: 12px;
  box-shadow: 0 6px 14px var(--primary-glow);
}

.cgpt-msg__user-icon {
  width: 34px;
  height: 34px;
  border-radius: 11px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 13px;
  font-weight: 700;
  color: #fff;
}

.cgpt-msg__content {
  flex: 1;
  min-width: 0;
  padding-top: 4px;
}

.cgpt-msg--user .cgpt-msg__content {
  padding: 12px 16px;
  border-radius: 16px;
  background: var(--surface);
  border: 1px solid var(--border);
  box-shadow: var(--shadow-xs);
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
  border-radius: 9px;
  background: transparent;
  color: var(--text-3);
  cursor: pointer;
  transition: all var(--t-fast);
}

.cgpt-msg-action:hover {
  background: var(--primary-soft);
  color: var(--primary);
}

.cgpt-dots span {
  display: inline-block;
  width: 8px;
  height: 8px;
  margin-right: 4px;
  background: var(--primary);
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
}

.cgpt-input-box {
  max-width: 820px;
  margin: 0 auto;
}

.cgpt-input-shell {
  display: flex;
  align-items: flex-end;
  gap: 8px;
  padding: 10px 10px 10px 18px;
  border: 1px solid var(--border-strong);
  border-radius: 24px;
  background: var(--surface);
  box-shadow: var(--shadow);
  transition: border-color var(--t-fast), box-shadow var(--t-fast);
}

.cgpt-input-shell--focus {
  border-color: var(--primary);
  box-shadow: 0 0 0 4px var(--primary-soft), var(--shadow);
}

.cgpt-input {
  flex: 1;
  border: none;
  outline: none;
  resize: none;
  font-size: 15px;
  line-height: 1.6;
  max-height: 200px;
  padding: 6px 0;
  font-family: inherit;
  background: transparent;
  color: var(--text);
}

.cgpt-input::placeholder {
  color: var(--text-3);
}

.cgpt-input-actions {
  flex-shrink: 0;
}

.cgpt-send {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 38px;
  height: 38px;
  border: none;
  border-radius: 50%;
  background: var(--surface-3);
  color: var(--text-3);
  cursor: not-allowed;
  transition: all var(--t-fast);
}

.cgpt-send--active {
  background: linear-gradient(135deg, var(--primary) 0%, var(--primary-strong) 100%);
  color: #fff;
  cursor: pointer;
  box-shadow: 0 6px 14px var(--primary-glow);
}

.cgpt-send--active:hover {
  transform: translateY(-1px);
}

.cgpt-send--stop {
  background: var(--ink);
  color: #fff;
  cursor: pointer;
}

.cgpt-disclaimer {
  margin: 10px 0 0;
  text-align: center;
  font-size: 12px;
  color: var(--text-3);
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
