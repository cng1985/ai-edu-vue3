<template>
  <div class="page">
    <PageHeader eyebrow="客户服务" title="客户咨询" subtitle="实时接收学员工单并在线回复，状态变更即时同步。">
      <template #title-extra>
        <span class="conn" :class="{ 'conn--on': wsConnected }">
          <i />{{ wsConnected ? '实时连接' : '连接中…' }}
        </span>
      </template>
      <div v-if="stats" class="stats">
        <span class="stats__item"><b class="num">{{ stats.total }}</b>全部</span>
        <span class="stats__item stats__item--warning"><b class="num">{{ stats.open }}</b>待处理</span>
        <span class="stats__item stats__item--primary"><b class="num">{{ stats.pending }}</b>处理中</span>
        <span class="stats__item stats__item--muted"><b class="num">{{ stats.closed }}</b>已关闭</span>
      </div>
    </PageHeader>

    <div class="support">
      <aside class="panel support__list">
        <div class="support__filter">
          <el-input v-model="filters.keyword" placeholder="搜索用户 / 主题" clearable :prefix-icon="Search" @keyup.enter="loadTickets" @clear="loadTickets" />
          <el-select v-model="filters.status" placeholder="状态" clearable style="width: 110px" @change="loadTickets">
            <el-option label="待处理" value="open" />
            <el-option label="处理中" value="pending" />
            <el-option label="已关闭" value="closed" />
          </el-select>
        </div>

        <div v-loading="loadingTickets" class="tickets">
          <button
            v-for="t in tickets"
            :key="t.id"
            type="button"
            class="ticket"
            :class="{ 'ticket--active': activeId === t.id }"
            @click="selectTicket(t)"
          >
            <span class="ticket__avatar" :style="{ background: avatarColor(t.userUsername) }">{{ (t.userNickname || t.userUsername || '?').slice(0, 1) }}</span>
            <span class="ticket__main">
              <span class="ticket__top">
                <strong>{{ t.subject }}</strong>
                <time>{{ shortTime(t.lastMessageAt) }}</time>
              </span>
              <span class="ticket__user">{{ t.userNickname }} · @{{ t.userUsername }}</span>
              <span class="ticket__preview">{{ t.lastMessage || '暂无消息' }}</span>
            </span>
            <span class="ticket__status" :class="`ticket__status--${t.status}`">{{ statusLabel(t.status) }}</span>
          </button>
          <el-empty v-if="!loadingTickets && !tickets.length" description="暂无咨询工单" :image-size="80" />
        </div>
      </aside>

      <section class="panel support__chat">
        <template v-if="activeTicket">
          <div class="chat-head">
            <span class="ticket__avatar ticket__avatar--lg" :style="{ background: avatarColor(activeTicket.userUsername) }">
              {{ (activeTicket.userNickname || activeTicket.userUsername || '?').slice(0, 1) }}
            </span>
            <div class="chat-head__text">
              <h3>{{ activeTicket.subject }}</h3>
              <span>{{ activeTicket.userNickname }} · @{{ activeTicket.userUsername }}</span>
            </div>
            <el-select
              v-if="auth.hasPermission(PERM.CUSTOMER_REPLY)"
              v-model="activeTicket.status"
              style="width: 120px"
              @change="updateStatus"
            >
              <el-option label="待处理" value="open" />
              <el-option label="处理中" value="pending" />
              <el-option label="已关闭" value="closed" />
            </el-select>
          </div>

          <div ref="scrollArea" v-loading="loadingMessages" class="messages">
            <div
              v-for="msg in messages"
              :key="msg.id"
              class="msg"
              :class="{ 'msg--staff': msg.senderRole !== 'learner' }"
            >
              <div class="msg__meta">
                <strong>{{ msg.senderRole === 'learner' ? (msg.senderNickname || '客户') : (msg.senderNickname || '客服') }}</strong>
                <time>{{ formatTime(msg.createdAt) }}</time>
              </div>
              <div class="msg__bubble">{{ msg.content }}</div>
            </div>
            <div v-if="!loadingMessages && !messages.length" class="messages__empty muted">还没有消息，先发一条打个招呼吧。</div>
          </div>

          <div v-if="activeTicket.status !== 'closed'" class="composer">
            <el-input
              v-model="replyText"
              type="textarea"
              :autosize="{ minRows: 1, maxRows: 5 }"
              resize="none"
              placeholder="输入回复内容，Enter 发送，Shift + Enter 换行"
              @keydown.enter.exact.prevent="sendReply"
            />
            <el-button
              v-permission="PERM.CUSTOMER_REPLY"
              type="primary"
              :loading="sending"
              :disabled="!replyText.trim()"
              class="composer__send"
              @click="sendReply"
            >
              <el-icon :size="16"><Promotion /></el-icon>
            </el-button>
          </div>
          <div v-else class="closed">
            <el-icon :size="16"><Lock /></el-icon>此工单已关闭，如需继续沟通请将状态改为「处理中」。
          </div>
        </template>
        <div v-else class="support__empty">
          <div class="support__empty-icon"><el-icon :size="28"><ChatDotRound /></el-icon></div>
          <h3>选择一个工单开始处理</h3>
          <p class="muted">左侧列表会实时显示新工单与最新消息。</p>
        </div>
      </section>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, onUnmounted, nextTick } from 'vue'
import { ElMessage } from 'element-plus'
import { Search, Promotion, Lock, ChatDotRound } from '@element-plus/icons-vue'
import { customersApi } from '../api'
import { useAuthStore } from '../stores/auth'
import { PERM } from '../constants/permissions'
import { useAdminWebSocket } from '../composables/useWebSocket'
import PageHeader from '../components/common/PageHeader.vue'

const AVATAR_COLORS = ['#6b5cff', '#0fb981', '#f59e0b', '#2a8cf4', '#f0589a', '#0ea5e9']
function avatarColor(seed = '') {
  let h = 0
  for (const ch of seed) h = (h * 31 + ch.charCodeAt(0)) >>> 0
  return AVATAR_COLORS[h % AVATAR_COLORS.length]
}

function shortTime(ts) {
  if (!ts) return ''
  const d = new Date(ts)
  const now = new Date()
  if (d.toDateString() === now.toDateString()) {
    return d.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })
  }
  return d.toLocaleDateString('zh-CN', { month: 'numeric', day: 'numeric' })
}

const auth = useAuthStore()
const ws = useAdminWebSocket()
const wsConnected = ws.connected

const stats = ref(null)
const tickets = ref([])
const messages = ref([])
const activeId = ref(null)
const activeTicket = ref(null)
const loadingTickets = ref(false)
const loadingMessages = ref(false)
const sending = ref(false)
const replyText = ref('')
const scrollArea = ref(null)
const filters = reactive({ keyword: '', status: '' })

let unsubscribers = []

const statusLabel = (s) => ({ open: '待处理', pending: '处理中', closed: '已关闭' }[s] || s)
const statusTagType = (s) => ({ open: 'warning', pending: '', closed: 'info' }[s] || '')

function formatTime(ts) {
  if (!ts) return ''
  return new Date(ts).toLocaleString('zh-CN')
}

async function loadStats() {
  try {
    stats.value = await customersApi.stats()
  } catch { /* ignore */ }
}

async function loadTickets() {
  loadingTickets.value = true
  try {
    const res = await customersApi.listTickets(filters)
    tickets.value = res.list || []
  } finally {
    loadingTickets.value = false
  }
}

async function selectTicket(ticket) {
  activeId.value = ticket.id
  activeTicket.value = ticket
  loadingMessages.value = true
  try {
    const res = await customersApi.listMessages(ticket.id)
    messages.value = res.list || []
    ws.send('support.subscribe', { ticketId: ticket.id })
    await nextTick()
    if (scrollArea.value) scrollArea.value.scrollTop = scrollArea.value.scrollHeight
  } finally {
    loadingMessages.value = false
  }
}

async function sendReply() {
  const content = replyText.value.trim()
  if (!content || !activeId.value) return
  sending.value = true
  try {
    const sent = ws.send('support.send', { ticketId: activeId.value, content })
    if (!sent) {
      const msg = await customersApi.reply(activeId.value, content)
      messages.value.push(msg)
    }
    replyText.value = ''
    await nextTick()
    if (scrollArea.value) scrollArea.value.scrollTop = scrollArea.value.scrollHeight
    loadStats()
  } catch (e) {
    ElMessage.error(e.message || '发送失败')
  } finally {
    sending.value = false
  }
}

async function updateStatus(status) {
  try {
    const updated = await customersApi.updateStatus(activeId.value, status)
    activeTicket.value = { ...activeTicket.value, ...updated }
    const idx = tickets.value.findIndex((t) => t.id === activeId.value)
    if (idx >= 0) tickets.value[idx] = { ...tickets.value[idx], ...updated }
    loadStats()
    ElMessage.success('状态已更新')
  } catch (e) {
    ElMessage.error(e.message || '更新失败')
  }
}

function setupWS() {
  ws.connect()
  unsubscribers.push(
    ws.on('support.message', (msg) => {
      if (msg.ticketId === activeId.value) {
        if (!messages.value.find((m) => m.id === msg.id)) {
          messages.value.push(msg)
          nextTick(() => {
            if (scrollArea.value) scrollArea.value.scrollTop = scrollArea.value.scrollHeight
          })
        }
      }
      const t = tickets.value.find((x) => x.id === msg.ticketId)
      if (t) {
        t.lastMessage = msg.content
        t.lastMessageAt = msg.createdAt
      }
    }),
    ws.on('support.ticket.new', (ticket) => {
      if (!tickets.value.find((t) => t.id === ticket.id)) {
        tickets.value.unshift(ticket)
        loadStats()
      }
    }),
    ws.on('support.ticket.update', (ticket) => {
      const idx = tickets.value.findIndex((t) => t.id === ticket.id)
      if (idx >= 0) tickets.value[idx] = { ...tickets.value[idx], ...ticket }
      if (activeId.value === ticket.id) {
        activeTicket.value = { ...activeTicket.value, ...ticket }
      }
      loadStats()
    })
  )
}

onMounted(() => {
  setupWS()
  loadStats()
  loadTickets()
})

onUnmounted(() => {
  unsubscribers.forEach((fn) => fn())
  ws.disconnect()
})
</script>

<style scoped>
.conn {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 3px 10px;
  border-radius: 999px;
  background: var(--surface-3);
  color: var(--text-3);
  font-size: 12px;
  font-weight: 600;
}

.conn i {
  width: 7px;
  height: 7px;
  border-radius: 50%;
  background: var(--text-3);
}

.conn--on {
  background: var(--success-soft);
  color: var(--success-strong);
}

.conn--on i {
  background: var(--success);
  box-shadow: 0 0 0 3px rgba(15, 185, 129, 0.2);
  animation: pulse 1.8s infinite;
}

@keyframes pulse {
  0%, 100% { box-shadow: 0 0 0 3px rgba(15, 185, 129, 0.2); }
  50% { box-shadow: 0 0 0 6px rgba(15, 185, 129, 0.08); }
}

.stats {
  display: flex;
  gap: 8px;
}

.stats__item {
  --tone: var(--text);
  --tone-soft: var(--surface-3);
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 6px 12px;
  border-radius: 999px;
  background: var(--tone-soft);
  color: var(--tone);
  font-size: 12.5px;
  font-weight: 600;
}

.stats__item b {
  font-size: 14px;
}

.stats__item--warning { --tone: var(--warning-strong); --tone-soft: var(--warning-soft); }
.stats__item--primary { --tone: var(--primary-deep); --tone-soft: var(--primary-soft); }
.stats__item--muted { --tone: var(--text-3); --tone-soft: var(--surface-3); }

.support {
  display: grid;
  grid-template-columns: 360px minmax(0, 1fr);
  gap: 20px;
  height: calc(100vh - var(--header-height) - 190px);
  min-height: 560px;
}

.support__list {
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.support__filter {
  display: flex;
  gap: 8px;
  padding: 14px;
  border-bottom: 1px solid var(--border);
}

.tickets {
  flex: 1;
  overflow-y: auto;
  padding: 8px;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.ticket {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  width: 100%;
  padding: 12px;
  border: 1px solid transparent;
  border-radius: 14px;
  background: transparent;
  text-align: left;
  cursor: pointer;
  transition: all var(--t-fast);
}

.ticket:hover {
  background: var(--surface-2);
}

.ticket--active {
  background: var(--primary-soft);
  border-color: var(--primary-soft-2);
}

.ticket__avatar {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 38px;
  height: 38px;
  flex-shrink: 0;
  border-radius: 12px;
  color: #fff;
  font-weight: 700;
  font-size: 14px;
}

.ticket__avatar--lg {
  width: 44px;
  height: 44px;
  border-radius: 14px;
  font-size: 16px;
}

.ticket__main {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.ticket__top {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
}

.ticket__top strong {
  font-size: 14px;
  color: var(--text);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.ticket__top time {
  font-size: 11.5px;
  color: var(--text-3);
  flex-shrink: 0;
}

.ticket__user {
  font-size: 12px;
  color: var(--text-3);
}

.ticket__preview {
  font-size: 12.5px;
  color: var(--text-2);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.ticket__status {
  flex-shrink: 0;
  align-self: flex-end;
  padding: 2px 8px;
  border-radius: 999px;
  font-size: 11px;
  font-weight: 700;
  background: var(--surface-3);
  color: var(--text-3);
}

.ticket__status--open { background: var(--warning-soft); color: var(--warning-strong); }
.ticket__status--pending { background: var(--primary-soft); color: var(--primary-deep); }

.support__chat {
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.chat-head {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 16px 20px;
  border-bottom: 1px solid var(--border);
}

.chat-head__text {
  flex: 1;
  min-width: 0;
  line-height: 1.3;
}

.chat-head__text h3 {
  margin: 0;
  font-size: 16px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.chat-head__text span {
  font-size: 12.5px;
  color: var(--text-3);
}

.messages {
  flex: 1;
  overflow-y: auto;
  padding: 20px;
  display: flex;
  flex-direction: column;
  gap: 16px;
  background:
    radial-gradient(500px 240px at 100% 0%, rgba(107, 92, 255, 0.05), transparent 60%),
    var(--surface-2);
}

.messages__empty {
  margin: auto;
  font-size: 13px;
}

.msg {
  max-width: 72%;
  align-self: flex-start;
}

.msg--staff {
  align-self: flex-end;
}

.msg__meta {
  display: flex;
  gap: 8px;
  font-size: 12px;
  color: var(--text-3);
  margin-bottom: 5px;
}

.msg--staff .msg__meta {
  justify-content: flex-end;
}

.msg__bubble {
  padding: 10px 14px;
  border-radius: 16px 16px 16px 6px;
  background: var(--surface);
  border: 1px solid var(--border);
  color: var(--text);
  font-size: 14px;
  line-height: 1.6;
  white-space: pre-wrap;
  word-break: break-word;
  box-shadow: var(--shadow-xs);
}

.msg--staff .msg__bubble {
  border-radius: 16px 16px 6px 16px;
  background: linear-gradient(135deg, var(--primary) 0%, var(--primary-strong) 100%);
  border-color: transparent;
  color: #fff;
  box-shadow: 0 6px 16px var(--primary-glow);
}

.composer {
  display: flex;
  gap: 10px;
  align-items: flex-end;
  padding: 14px 16px;
  border-top: 1px solid var(--border);
  background: var(--surface);
}

.composer .el-textarea {
  flex: 1;
}

.composer :deep(.el-textarea__inner) {
  padding: 10px 14px;
  border-radius: 14px !important;
}

.composer__send {
  width: 44px;
  height: 44px;
  padding: 0;
  border-radius: 13px;
  flex-shrink: 0;
}

.closed {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 14px 20px;
  border-top: 1px solid var(--border);
  background: var(--surface-2);
  color: var(--text-3);
  font-size: 13px;
}

.support__empty {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  text-align: center;
  padding: 40px;
}

.support__empty-icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 64px;
  height: 64px;
  border-radius: 20px;
  background: var(--primary-soft);
  color: var(--primary);
  margin-bottom: 16px;
}

.support__empty h3 {
  margin: 0;
  font-size: 17px;
}

.support__empty p {
  margin: 6px 0 0;
  font-size: 13.5px;
}

@media (max-width: 1000px) {
  .support {
    grid-template-columns: 1fr;
    height: auto;
  }

  .support__list {
    max-height: 420px;
  }

  .support__chat {
    min-height: 520px;
  }
}
</style>
