<script setup>
import { ref, nextTick, watch, onMounted, onUnmounted } from 'vue'
import { useCustomerStore } from '../stores/customer'

const customer = useCustomerStore()
const input = ref('')
const newSubject = ref('')
const showNewForm = ref(false)
const scrollArea = ref(null)

const statusLabel = { open: '待处理', pending: '处理中', closed: '已关闭' }
const statusClass = { open: 'status--open', pending: 'status--pending', closed: 'status--closed' }

function formatTime(ts) {
  if (!ts) return ''
  return new Date(ts).toLocaleString('zh-CN', { month: '2-digit', day: '2-digit', hour: '2-digit', minute: '2-digit' })
}

function isMine(msg) {
  return msg.senderRole === 'learner'
}

async function scrollToBottom() {
  await nextTick()
  if (scrollArea.value) scrollArea.value.scrollTop = scrollArea.value.scrollHeight
}

async function handleSend() {
  const text = input.value.trim()
  if (!text) return
  await customer.sendMessage(text)
  input.value = ''
  scrollToBottom()
}

async function handleCreate() {
  const content = input.value.trim()
  if (!content) return
  await customer.createTicket(newSubject.value.trim() || '咨询求助', content)
  showNewForm.value = false
  newSubject.value = ''
  input.value = ''
  scrollToBottom()
}

function onKeydown(e) {
  if (e.key === 'Enter' && !e.shiftKey && !e.isComposing) {
    e.preventDefault()
    if (showNewForm.value || !customer.activeTicketId) handleCreate()
    else handleSend()
  }
}

watch(() => customer.messages.length, scrollToBottom)

onMounted(async () => {
  customer.setupWS()
  await customer.loadTickets()
  if (customer.tickets.length > 0) {
    await customer.selectTicket(customer.tickets[0].id)
  } else {
    showNewForm.value = true
  }
  scrollToBottom()
})

onUnmounted(() => customer.teardownWS())
</script>

<template>
  <div class="support">
    <header class="support__header">
      <div>
        <h1>客户咨询</h1>
        <p>
          在线联系客服，实时沟通学习问题与平台使用疑问。
          <span class="support__conn" :class="{ 'support__conn--on': customer.connected }">
            {{ customer.connected ? '● 已连接' : '○ 连接中…' }}
          </span>
        </p>
      </div>
      <button class="btn btn--primary" @click="showNewForm = true; customer.activeTicketId = null">
        新建咨询
      </button>
    </header>

    <div class="support__body">
      <aside class="support__sidebar">
        <div v-if="customer.loading && !customer.tickets.length" class="support__empty-list">加载中…</div>
        <div v-else-if="!customer.tickets.length" class="support__empty-list">暂无咨询记录</div>
        <button
          v-for="t in customer.tickets"
          :key="t.id"
          class="support__ticket"
          :class="{ 'support__ticket--active': t.id === customer.activeTicketId }"
          @click="showNewForm = false; customer.selectTicket(t.id)"
        >
          <div class="support__ticket-top">
            <strong>{{ t.subject }}</strong>
            <span class="support__status" :class="statusClass[t.status]">{{ statusLabel[t.status] }}</span>
          </div>
          <p class="support__ticket-preview">{{ t.lastMessage || '暂无消息' }}</p>
          <time>{{ formatTime(t.lastMessageAt) }}</time>
        </button>
      </aside>

      <main class="support__chat">
        <div v-if="showNewForm || !customer.activeTicketId" class="support__new">
          <h2>发起新咨询</h2>
          <input v-model="newSubject" class="support__input-subject" placeholder="咨询主题（可选）" />
          <textarea
            v-model="input"
            class="support__textarea"
            rows="5"
            placeholder="请描述您的问题…"
            @keydown="onKeydown"
          />
          <button class="btn btn--primary" :disabled="customer.sending || !input.trim()" @click="handleCreate">
            {{ customer.sending ? '提交中…' : '提交咨询' }}
          </button>
        </div>

        <template v-else>
          <div class="support__chat-header">
            <h2>{{ customer.activeTicket?.subject }}</h2>
            <span class="support__status" :class="statusClass[customer.activeTicket?.status]">
              {{ statusLabel[customer.activeTicket?.status] }}
            </span>
          </div>

          <div ref="scrollArea" class="support__messages">
            <div
              v-for="msg in customer.messages"
              :key="msg.id"
              class="support__msg"
              :class="{ 'support__msg--mine': isMine(msg) }"
            >
              <div class="support__msg-meta">
                <strong>{{ isMine(msg) ? '我' : (msg.senderNickname || '客服') }}</strong>
                <time>{{ formatTime(msg.createdAt) }}</time>
              </div>
              <div class="support__msg-bubble">{{ msg.content }}</div>
            </div>
            <div v-if="customer.activeTicket?.status === 'closed'" class="support__closed-tip">
              此咨询已关闭，如需帮助请新建咨询
            </div>
          </div>

          <div v-if="customer.activeTicket?.status !== 'closed'" class="support__composer">
            <textarea
              v-model="input"
              class="support__textarea"
              rows="2"
              placeholder="输入消息，Enter 发送…"
              @keydown="onKeydown"
            />
            <button class="btn btn--primary" :disabled="customer.sending || !input.trim()" @click="handleSend">
              发送
            </button>
          </div>
        </template>
      </main>
    </div>
  </div>
</template>

<style scoped>
.support {
  max-width: 1100px;
  margin: 0 auto;
  padding: 0 4px;
}

.support__header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 16px;
  margin-bottom: 20px;
}

.support__header h1 {
  font-size: 24px;
  font-weight: 700;
  margin-bottom: 6px;
}

.support__header p {
  color: var(--text-2);
  font-size: 14px;
}

.support__conn {
  margin-left: 8px;
  font-size: 12px;
  color: var(--text-3);
}

.support__conn--on {
  color: #10b981;
}

.support__body {
  display: grid;
  grid-template-columns: 280px 1fr;
  gap: 16px;
  min-height: 520px;
}

.support__sidebar {
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: var(--radius);
  overflow-y: auto;
  max-height: 600px;
}

.support__empty-list {
  padding: 24px;
  text-align: center;
  color: var(--text-3);
  font-size: 14px;
}

.support__ticket {
  display: block;
  width: 100%;
  text-align: left;
  padding: 14px 16px;
  border: none;
  border-bottom: 1px solid var(--border);
  background: transparent;
  cursor: pointer;
  transition: background 0.15s;
}

.support__ticket:hover {
  background: var(--surface-2);
}

.support__ticket--active {
  background: var(--primary-soft);
}

.support__ticket-top {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 8px;
  margin-bottom: 4px;
}

.support__ticket-top strong {
  font-size: 14px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.support__ticket-preview {
  font-size: 12.5px;
  color: var(--text-3);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  margin-bottom: 4px;
}

.support__ticket time {
  font-size: 11px;
  color: var(--text-3);
}

.support__chat {
  display: flex;
  flex-direction: column;
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: var(--radius);
  min-height: 520px;
}

.support__chat-header {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 16px 20px;
  border-bottom: 1px solid var(--border);
}

.support__chat-header h2 {
  font-size: 16px;
  font-weight: 600;
}

.support__status {
  font-size: 11px;
  font-weight: 600;
  padding: 2px 8px;
  border-radius: 999px;
}

.status--open { background: #fef3c7; color: #b45309; }
.status--pending { background: #dbeafe; color: #1d4ed8; }
.status--closed { background: #f3f4f6; color: #6b7280; }

.support__messages {
  flex: 1;
  overflow-y: auto;
  padding: 20px;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.support__msg {
  max-width: 75%;
  align-self: flex-start;
}

.support__msg--mine {
  align-self: flex-end;
}

.support__msg-meta {
  display: flex;
  gap: 8px;
  align-items: center;
  margin-bottom: 4px;
  font-size: 12px;
  color: var(--text-3);
}

.support__msg--mine .support__msg-meta {
  justify-content: flex-end;
}

.support__msg-bubble {
  padding: 10px 14px;
  border-radius: 12px;
  background: var(--surface-2);
  font-size: 14px;
  line-height: 1.55;
  white-space: pre-wrap;
  word-break: break-word;
}

.support__msg--mine .support__msg-bubble {
  background: var(--primary);
  color: #fff;
}

.support__closed-tip {
  text-align: center;
  color: var(--text-3);
  font-size: 13px;
  padding: 12px;
}

.support__composer,
.support__new {
  padding: 16px 20px;
  border-top: 1px solid var(--border);
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.support__new h2 {
  font-size: 16px;
  font-weight: 600;
}

.support__input-subject {
  padding: 10px 12px;
  border: 1px solid var(--border);
  border-radius: var(--radius-sm);
  font-size: 14px;
  background: var(--surface-2);
}

.support__textarea {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid var(--border);
  border-radius: var(--radius-sm);
  font-size: 14px;
  resize: vertical;
  background: var(--surface-2);
  font-family: inherit;
}

.support__composer {
  flex-direction: row;
  align-items: flex-end;
}

.support__composer .support__textarea {
  flex: 1;
}

@media (max-width: 768px) {
  .support__body {
    grid-template-columns: 1fr;
  }
  .support__sidebar {
    max-height: 200px;
  }
}
</style>
