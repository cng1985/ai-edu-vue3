<script setup>
import { ref, nextTick, watch, onMounted, onUnmounted } from 'vue'
import { useCustomerStore } from '../stores/customer'
import Icon from '../components/Icon.vue'

const customer = useCustomerStore()
const input = ref('')
const newSubject = ref('')
const showNewForm = ref(false)
const scrollArea = ref(null)

const statusLabel = { open: '待处理', pending: '处理中', closed: '已关闭' }
const statusClass = { open: 'tag--warning', pending: 'tag--info', closed: 'tag--neutral' }

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

function startNew() {
  showNewForm.value = true
  customer.activeTicketId = null
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
  <div class="page support">
    <header class="page-header page-header--split">
      <div>
        <span class="eyebrow">在线服务</span>
        <h1>客户咨询</h1>
        <p>
          在线联系客服，实时沟通学习问题与平台使用疑问。
        </p>
      </div>
      <div class="support__head-actions">
        <span class="support__conn" :class="{ 'support__conn--on': customer.connected }">
          <i></i>{{ customer.connected ? '已连接' : '连接中…' }}
        </span>
        <button class="btn btn--primary" @click="startNew">
          <Icon name="plus" :size="15" :stroke="2.5" /> 新建咨询
        </button>
      </div>
    </header>

    <div class="support__body">
      <aside class="support__sidebar card">
        <div class="support__sidebar-head">
          <strong>我的咨询</strong>
          <span class="tag tag--neutral">{{ customer.tickets.length }}</span>
        </div>
        <div v-if="customer.loading && !customer.tickets.length" class="support__empty-list">加载中…</div>
        <div v-else-if="!customer.tickets.length" class="support__empty-list">
          <Icon name="headset" :size="28" />
          <span>暂无咨询记录</span>
        </div>
        <button
          v-for="t in customer.tickets"
          :key="t.id"
          class="support__ticket"
          :class="{ 'support__ticket--active': t.id === customer.activeTicketId && !showNewForm }"
          @click="showNewForm = false; customer.selectTicket(t.id)"
        >
          <div class="support__ticket-top">
            <strong>{{ t.subject }}</strong>
            <span class="tag" :class="statusClass[t.status]">{{ statusLabel[t.status] }}</span>
          </div>
          <p class="support__ticket-preview">{{ t.lastMessage || '暂无消息' }}</p>
          <time>{{ formatTime(t.lastMessageAt) }}</time>
        </button>
      </aside>

      <main class="support__chat card">
        <div v-if="showNewForm || !customer.activeTicketId" class="support__new">
          <div class="support__new-icon"><Icon name="headset" :size="26" /></div>
          <h2>发起新咨询</h2>
          <p>描述你遇到的问题，客服会尽快回复。</p>
          <input v-model="newSubject" class="input" placeholder="咨询主题（可选）" />
          <textarea
            v-model="input"
            class="textarea"
            rows="5"
            placeholder="请描述您的问题…"
            @keydown="onKeydown"
          />
          <button class="btn btn--primary btn--lg" :disabled="customer.sending || !input.trim()" @click="handleCreate">
            <Icon name="send" :size="15" /> {{ customer.sending ? '提交中…' : '提交咨询' }}
          </button>
        </div>

        <template v-else>
          <div class="support__chat-header">
            <div>
              <h2>{{ customer.activeTicket?.subject }}</h2>
              <small>工单 #{{ String(customer.activeTicket?.id || '').slice(-6) }}</small>
            </div>
            <span class="tag" :class="statusClass[customer.activeTicket?.status]">
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
              <Icon name="lock" :size="13" /> 此咨询已关闭，如需帮助请新建咨询
            </div>
          </div>

          <div v-if="customer.activeTicket?.status !== 'closed'" class="support__composer">
            <textarea
              v-model="input"
              rows="2"
              placeholder="输入消息，Enter 发送…"
              @keydown="onKeydown"
            />
            <button class="support__send" :disabled="customer.sending || !input.trim()" title="发送" @click="handleSend">
              <Icon name="send" :size="16" />
            </button>
          </div>
        </template>
      </main>
    </div>
  </div>
</template>

<style scoped>
.support {
  max-width: 1160px;
}

.support__head-actions {
  display: flex;
  align-items: center;
  gap: 12px;
}

.support__conn {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 6px 12px;
  border-radius: 999px;
  background: var(--surface-3);
  font-size: 12.5px;
  font-weight: 600;
  color: var(--text-3);
}

.support__conn i {
  width: 7px;
  height: 7px;
  border-radius: 50%;
  background: var(--text-3);
}

.support__conn--on {
  background: var(--success-soft);
  color: var(--success-strong);
}

.support__conn--on i {
  background: var(--success);
  box-shadow: 0 0 0 3px rgba(15, 185, 129, 0.2);
}

.support__body {
  display: grid;
  grid-template-columns: 300px 1fr;
  gap: 18px;
  min-height: 560px;
}

.support__sidebar {
  display: flex;
  flex-direction: column;
  overflow: hidden;
  max-height: 640px;
}

.support__sidebar-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 18px 12px;
  border-bottom: 1px solid var(--border);
  font-size: 14px;
}

.support__empty-list {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 10px;
  padding: 40px 24px;
  text-align: center;
  color: var(--text-3);
  font-size: 13.5px;
}

.support__ticket {
  display: block;
  width: 100%;
  text-align: left;
  padding: 14px 18px;
  border: none;
  border-left: 3px solid transparent;
  border-bottom: 1px solid var(--border);
  background: transparent;
  font: inherit;
  cursor: pointer;
  transition: background var(--t-fast);
}

.support__ticket:hover {
  background: var(--surface-2);
}

.support__ticket--active {
  background: var(--primary-soft);
  border-left-color: var(--primary);
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
  margin: 0 0 6px;
  font-size: 12.5px;
  color: var(--text-2);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.support__ticket time {
  font-size: 11.5px;
  color: var(--text-3);
}

.support__chat {
  display: flex;
  flex-direction: column;
  min-height: 560px;
  overflow: hidden;
}

.support__chat-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
  padding: 16px 22px;
  border-bottom: 1px solid var(--border);
}

.support__chat-header h2 {
  margin: 0;
  font-size: 16px;
}

.support__chat-header small {
  color: var(--text-3);
  font-size: 12px;
}

.support__messages {
  flex: 1;
  overflow-y: auto;
  padding: 22px;
  display: flex;
  flex-direction: column;
  gap: 16px;
  background: var(--surface-2);
}

.support__msg {
  max-width: 72%;
  align-self: flex-start;
  animation: fade-up 0.25s var(--ease) both;
}

.support__msg--mine {
  align-self: flex-end;
}

.support__msg-meta {
  display: flex;
  gap: 8px;
  align-items: center;
  margin-bottom: 5px;
  font-size: 12px;
  color: var(--text-3);
}

.support__msg-meta strong {
  color: var(--text-2);
}

.support__msg--mine .support__msg-meta {
  justify-content: flex-end;
}

.support__msg-bubble {
  padding: 11px 15px;
  border-radius: 16px;
  border-top-left-radius: 6px;
  background: var(--surface);
  border: 1px solid var(--border);
  font-size: 14px;
  line-height: 1.6;
  white-space: pre-wrap;
  word-break: break-word;
  box-shadow: var(--shadow-xs);
}

.support__msg--mine .support__msg-bubble {
  background: linear-gradient(135deg, var(--primary), var(--primary-strong));
  border-color: transparent;
  color: #fff;
  border-radius: 16px;
  border-top-right-radius: 6px;
  box-shadow: 0 6px 16px var(--primary-glow);
}

.support__closed-tip {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  color: var(--text-3);
  font-size: 13px;
  padding: 12px;
}

.support__new {
  display: flex;
  flex-direction: column;
  gap: 12px;
  max-width: 520px;
  margin: auto;
  padding: 40px 24px;
  width: 100%;
}

.support__new-icon {
  display: grid;
  place-items: center;
  width: 60px;
  height: 60px;
  border-radius: 20px;
  background: var(--primary-soft);
  color: var(--primary-strong);
}

.support__new h2 {
  margin: 6px 0 0;
  font-size: 22px;
}

.support__new p {
  margin: -6px 0 8px;
  color: var(--text-2);
  font-size: 14px;
}

.support__composer {
  display: flex;
  gap: 10px;
  align-items: flex-end;
  padding: 14px 18px;
  border-top: 1px solid var(--border);
}

.support__composer textarea {
  flex: 1;
  padding: 11px 14px;
  border: 1px solid var(--border-strong);
  border-radius: 14px;
  font: inherit;
  font-size: 14px;
  line-height: 1.55;
  resize: none;
  background: var(--surface);
  outline: none;
  transition: border-color var(--t-fast), box-shadow var(--t-fast);
}

.support__composer textarea:focus {
  border-color: var(--primary);
  box-shadow: 0 0 0 4px var(--primary-soft-2);
}

.support__send {
  display: grid;
  place-items: center;
  width: 44px;
  height: 44px;
  border: none;
  border-radius: 13px;
  background: linear-gradient(135deg, var(--primary), var(--primary-strong));
  color: #fff;
  cursor: pointer;
  box-shadow: 0 6px 14px var(--primary-glow);
  transition: transform var(--t-fast), opacity var(--t-fast);
}

.support__send:hover:not(:disabled) {
  transform: translateY(-1px);
}

.support__send:disabled {
  opacity: 0.4;
  box-shadow: none;
  cursor: not-allowed;
}

@media (max-width: 860px) {
  .support__body {
    grid-template-columns: 1fr;
  }

  .support__sidebar {
    max-height: 240px;
  }
}
</style>
