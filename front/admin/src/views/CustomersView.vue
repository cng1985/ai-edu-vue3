<template>
  <div>
    <div class="page-header">
      <div>
        <h2>客户咨询</h2>
        <span class="conn-badge" :class="{ 'conn-badge--on': wsConnected }">
          {{ wsConnected ? '● 实时连接' : '○ 连接中' }}
        </span>
      </div>
      <div v-if="stats" class="stats-row">
        <el-tag>全部 {{ stats.total }}</el-tag>
        <el-tag type="warning">待处理 {{ stats.open }}</el-tag>
        <el-tag type="primary">处理中 {{ stats.pending }}</el-tag>
        <el-tag type="info">已关闭 {{ stats.closed }}</el-tag>
      </div>
    </div>

    <div class="customer-layout">
      <el-card shadow="never" class="ticket-list-card">
        <div class="filter-bar">
          <el-input v-model="filters.keyword" placeholder="搜索用户/主题" clearable style="width: 180px" @keyup.enter="loadTickets" />
          <el-select v-model="filters.status" placeholder="状态" clearable style="width: 120px" @change="loadTickets">
            <el-option label="待处理" value="open" />
            <el-option label="处理中" value="pending" />
            <el-option label="已关闭" value="closed" />
          </el-select>
          <el-button type="primary" @click="loadTickets">查询</el-button>
        </div>

        <div v-loading="loadingTickets" class="ticket-list">
          <div
            v-for="t in tickets"
            :key="t.id"
            class="ticket-item"
            :class="{ 'ticket-item--active': activeId === t.id }"
            @click="selectTicket(t)"
          >
            <div class="ticket-item__top">
              <strong>{{ t.subject }}</strong>
              <el-tag :type="statusTagType(t.status)" size="small">{{ statusLabel(t.status) }}</el-tag>
            </div>
            <div class="ticket-item__user">{{ t.userNickname }} (@{{ t.userUsername }})</div>
            <div class="ticket-item__preview">{{ t.lastMessage || '暂无消息' }}</div>
            <div class="ticket-item__time">{{ formatTime(t.lastMessageAt) }}</div>
          </div>
          <el-empty v-if="!loadingTickets && !tickets.length" description="暂无咨询工单" />
        </div>
      </el-card>

      <el-card shadow="never" class="chat-card">
        <template v-if="activeTicket">
          <div class="chat-header">
            <div>
              <h3>{{ activeTicket.subject }}</h3>
              <span class="chat-user">{{ activeTicket.userNickname }} (@{{ activeTicket.userUsername }})</span>
            </div>
            <div class="chat-actions">
              <el-select
                v-if="auth.hasPermission(PERM.CUSTOMER_REPLY)"
                v-model="activeTicket.status"
                size="small"
                style="width: 110px"
                @change="updateStatus"
              >
                <el-option label="待处理" value="open" />
                <el-option label="处理中" value="pending" />
                <el-option label="已关闭" value="closed" />
              </el-select>
            </div>
          </div>

          <div ref="scrollArea" v-loading="loadingMessages" class="chat-messages">
            <div
              v-for="msg in messages"
              :key="msg.id"
              class="chat-msg"
              :class="{ 'chat-msg--staff': msg.senderRole !== 'learner' }"
            >
              <div class="chat-msg__meta">
                <strong>{{ msg.senderRole === 'learner' ? (msg.senderNickname || '客户') : (msg.senderNickname || '客服') }}</strong>
                <time>{{ formatTime(msg.createdAt) }}</time>
              </div>
              <div class="chat-msg__bubble">{{ msg.content }}</div>
            </div>
          </div>

          <div v-if="activeTicket.status !== 'closed'" class="chat-composer">
            <el-input
              v-model="replyText"
              type="textarea"
              :rows="2"
              placeholder="输入回复内容…"
              @keydown.enter.exact.prevent="sendReply"
            />
            <el-button
              v-permission="PERM.CUSTOMER_REPLY"
              type="primary"
              :loading="sending"
              :disabled="!replyText.trim()"
              @click="sendReply"
            >
              回复
            </el-button>
          </div>
          <el-alert v-else title="此工单已关闭" type="info" :closable="false" show-icon />
        </template>
        <el-empty v-else description="请选择左侧工单查看详情" />
      </el-card>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, onUnmounted, nextTick } from 'vue'
import { ElMessage } from 'element-plus'
import { customersApi } from '../api'
import { useAuthStore } from '../stores/auth'
import { PERM } from '../constants/permissions'
import { useAdminWebSocket } from '../composables/useWebSocket'

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
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.page-header h2 {
  display: inline;
  margin-right: 10px;
}

.conn-badge {
  font-size: 12px;
  color: #909399;
}

.conn-badge--on {
  color: #67c23a;
}

.stats-row {
  display: flex;
  gap: 8px;
}

.customer-layout {
  display: grid;
  grid-template-columns: 360px 1fr;
  gap: 16px;
  min-height: 560px;
}

.ticket-list {
  max-height: 520px;
  overflow-y: auto;
}

.filter-bar {
  display: flex;
  gap: 8px;
  margin-bottom: 12px;
}

.ticket-item {
  padding: 12px;
  border-bottom: 1px solid #ebeef5;
  cursor: pointer;
  transition: background 0.15s;
}

.ticket-item:hover {
  background: #f5f7fa;
}

.ticket-item--active {
  background: #ecf5ff;
}

.ticket-item__top {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 8px;
  margin-bottom: 4px;
}

.ticket-item__top strong {
  font-size: 14px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.ticket-item__user {
  font-size: 12px;
  color: #909399;
  margin-bottom: 4px;
}

.ticket-item__preview {
  font-size: 12px;
  color: #606266;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.ticket-item__time {
  font-size: 11px;
  color: #c0c4cc;
  margin-top: 4px;
}

.chat-card {
  display: flex;
  flex-direction: column;
}

.chat-card :deep(.el-card__body) {
  display: flex;
  flex-direction: column;
  flex: 1;
  min-height: 520px;
}

.chat-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  padding-bottom: 12px;
  border-bottom: 1px solid #ebeef5;
  margin-bottom: 12px;
}

.chat-header h3 {
  margin: 0 0 4px;
  font-size: 16px;
}

.chat-user {
  font-size: 12px;
  color: #909399;
}

.chat-messages {
  flex: 1;
  overflow-y: auto;
  padding: 8px 0;
  display: flex;
  flex-direction: column;
  gap: 12px;
  min-height: 300px;
}

.chat-msg {
  max-width: 75%;
  align-self: flex-start;
}

.chat-msg--staff {
  align-self: flex-end;
}

.chat-msg__meta {
  display: flex;
  gap: 8px;
  font-size: 12px;
  color: #909399;
  margin-bottom: 4px;
}

.chat-msg--staff .chat-msg__meta {
  justify-content: flex-end;
}

.chat-msg__bubble {
  padding: 8px 12px;
  border-radius: 8px;
  background: #f4f4f5;
  font-size: 14px;
  line-height: 1.5;
  white-space: pre-wrap;
  word-break: break-word;
}

.chat-msg--staff .chat-msg__bubble {
  background: #409eff;
  color: #fff;
}

.chat-composer {
  display: flex;
  gap: 8px;
  align-items: flex-end;
  padding-top: 12px;
  border-top: 1px solid #ebeef5;
}

.chat-composer .el-textarea {
  flex: 1;
}

@media (max-width: 900px) {
  .customer-layout {
    grid-template-columns: 1fr;
  }
}
</style>
