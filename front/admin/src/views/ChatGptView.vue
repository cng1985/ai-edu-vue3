<template>
  <div class="chatgpt-page">
    <aside class="chatgpt-sidebar">
      <div class="chatgpt-sidebar__brand">
        <span class="chatgpt-sidebar__logo">GPT</span>
        <span v-if="!collapsed" class="chatgpt-sidebar__title">ChatGPT</span>
      </div>
      <el-button class="chatgpt-new-btn" type="primary" plain @click="clearChat">
        <el-icon><Plus /></el-icon>
        <span v-if="!collapsed">新对话</span>
      </el-button>
      <div v-if="!collapsed" class="chatgpt-sidebar__hint">
        管理端智能对话，支持切换虚拟模型与纯对话 / 知识库模式。
      </div>
    </aside>

    <main class="chatgpt-main">
      <header class="chatgpt-header">
        <div class="chatgpt-header__left">
          <h2>ChatGPT</h2>
          <el-tag size="small" :type="aiConfig?.enabled ? 'success' : 'info'">
            {{ statusLabel }}
          </el-tag>
        </div>
        <div class="chatgpt-header__controls">
          <el-select v-model="virtualModel" placeholder="虚拟模型" style="width: 180px" @change="onModelChange">
            <el-option
              v-for="vm in virtualModels"
              :key="vm.code"
              :label="`${vm.name} (${vm.code})`"
              :value="vm.code"
            />
          </el-select>
          <el-select v-model="mode" style="width: 140px">
            <el-option label="纯对话" value="chat" />
            <el-option label="知识库增强" value="rag" />
          </el-select>
          <el-button v-if="generating" type="danger" plain @click="stop">停止</el-button>
        </div>
      </header>

      <div ref="scrollArea" class="chatgpt-scroll">
        <div v-if="messages.length === 0" class="chatgpt-empty">
          <div class="chatgpt-empty__icon">💬</div>
          <h3>有什么可以帮你的？</h3>
          <p>选择虚拟模型后开始对话。纯对话模式类似 ChatGPT；知识库模式会检索课程内容。</p>
          <div class="chatgpt-suggestions">
            <button v-for="s in suggestions" :key="s" class="chatgpt-suggestion" @click="send(s)">{{ s }}</button>
          </div>
        </div>

        <div
          v-for="msg in messages"
          :key="msg.id"
          class="chatgpt-row"
          :class="`chatgpt-row--${msg.role}`"
        >
          <div class="chatgpt-avatar" :class="{ 'chatgpt-avatar--user': msg.role === 'user' }">
            {{ msg.role === 'user' ? '我' : 'AI' }}
          </div>
          <div class="chatgpt-bubble">
            <div v-if="msg.streaming && !msg.text" class="chatgpt-thinking">
              <span></span><span></span><span></span>
            </div>
            <div v-else class="chatgpt-markdown" v-html="renderMarkdown(msg.text)"></div>
            <div v-if="msg.meta && !msg.streaming" class="chatgpt-meta">
              {{ msg.meta }}
            </div>
          </div>
        </div>
      </div>

      <footer class="chatgpt-footer">
        <div class="chatgpt-input-wrap">
          <el-input
            v-model="input"
            type="textarea"
            :autosize="{ minRows: 1, maxRows: 4 }"
            placeholder="输入消息，Enter 发送，Shift+Enter 换行"
            :disabled="generating"
            @keydown="onKeydown"
          />
          <el-button type="primary" :loading="generating" :disabled="!input.trim()" @click="send()">
            发送
          </el-button>
        </div>
      </footer>
    </main>
  </div>
</template>

<script setup>
import { computed, nextTick, onMounted, ref, watch } from 'vue'
import { marked } from 'marked'
import { ElMessage } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import { aiApi } from '../api/ai.js'
import { aiModelsApi } from '../api/aiModels.js'

const collapsed = ref(false)
const input = ref('')
const messages = ref([])
const generating = ref(false)
const scrollArea = ref(null)
const aiConfig = ref(null)
const virtualModels = ref([])
const virtualModel = ref('chat-default')
const mode = ref('chat')
let cancelStream = null

const suggestions = [
  '用三句话解释什么是 RAG',
  '写一段 Python 快速排序代码',
  '帮我润色一段产品介绍',
  '平台默认使用哪个虚拟模型？'
]

const statusLabel = computed(() => {
  if (!aiConfig.value) return '连接中…'
  if (!aiConfig.value.enabled) return '本地模式'
  const parts = [aiConfig.value.model]
  if (aiConfig.value.provider) parts.push(aiConfig.value.provider)
  return parts.join(' · ')
})

function renderMarkdown(text) {
  if (!text) return ''
  return marked.parse(text, { breaks: true })
}

function buildHistory() {
  return messages.value
    .filter((m) => !m.streaming && m.text)
    .slice(-16)
    .map((m) => ({ role: m.role, content: m.text }))
}

function scrollToBottom() {
  nextTick(() => {
    if (scrollArea.value) scrollArea.value.scrollTop = scrollArea.value.scrollHeight
  })
}

watch(() => messages.value.map((m) => m.text.length).join(','), scrollToBottom)

async function loadMeta() {
  try {
    aiConfig.value = await aiApi.config()
    if (aiConfig.value?.defaultVirtualModel) {
      virtualModel.value = aiConfig.value.defaultVirtualModel
    }
    virtualModels.value = await aiModelsApi.listVirtualModelOptions()
    if (!virtualModels.value.length) {
      virtualModels.value = [
        { code: 'chat-default', name: '默认对话模型' },
        { code: 'chat-smart', name: '高级对话模型' }
      ]
    }
  } catch {
    aiConfig.value = { enabled: false }
  }
}

function onModelChange() {
  ElMessage.info(`已切换虚拟模型：${virtualModel.value}`)
}

function send(text) {
  const question = (text ?? input.value).trim()
  if (!question || generating.value) return

  messages.value.push({
    id: `${Date.now()}:user`,
    role: 'user',
    text: question,
    streaming: false
  })

  const reply = {
    id: `${Date.now()}:assistant`,
    role: 'assistant',
    text: '',
    streaming: true,
    meta: ''
  }
  messages.value.push(reply)
  generating.value = true
  input.value = ''

  const history = buildHistory().slice(0, -1)

  cancelStream = aiApi.chatStream(
    question,
    history,
    { virtualModel: virtualModel.value, mode: mode.value },
    {
      onToken: (chunk) => {
        reply.text += chunk
      },
      onDone: (payload) => {
        reply.text = payload.text || reply.text
        reply.streaming = false
        if (payload.model || payload.virtualModel) {
          reply.meta = [
            payload.virtualModel && `虚拟模型 ${payload.virtualModel}`,
            payload.canonicalModel && `统一模型 ${payload.canonicalModel}`,
            payload.provider && `厂商 ${payload.provider}`,
            payload.model && `调用 ${payload.model}`
          ].filter(Boolean).join(' · ')
        }
        generating.value = false
        cancelStream = null
        loadMeta()
      },
      onError: (err) => {
        reply.text = `请求失败：${err.message}`
        reply.streaming = false
        generating.value = false
        cancelStream = null
      }
    }
  )
}

function stop() {
  if (cancelStream) {
    cancelStream()
    cancelStream = null
  }
  const last = messages.value[messages.value.length - 1]
  if (last?.streaming) {
    last.streaming = false
    last.text += '\n\n*（已停止生成）*'
  }
  generating.value = false
}

function clearChat() {
  stop()
  messages.value = []
}

function onKeydown(e) {
  if (e.key === 'Enter' && !e.shiftKey && !e.isComposing) {
    e.preventDefault()
    send()
  }
}

onMounted(loadMeta)
</script>

<style scoped>
.chatgpt-page {
  display: flex;
  height: calc(100vh - 56px);
  margin: -20px;
  background: #f8fafc;
}
.chatgpt-sidebar {
  width: 220px;
  background: #111827;
  color: #e5e7eb;
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 12px;
}
.chatgpt-sidebar__brand {
  display: flex;
  align-items: center;
  gap: 10px;
}
.chatgpt-sidebar__logo {
  width: 32px;
  height: 32px;
  border-radius: 8px;
  background: linear-gradient(135deg, #10a37f, #1a7f64);
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 700;
  font-size: 12px;
}
.chatgpt-sidebar__title {
  font-weight: 600;
}
.chatgpt-new-btn {
  width: 100%;
}
.chatgpt-sidebar__hint {
  font-size: 12px;
  color: #9ca3af;
  line-height: 1.5;
}
.chatgpt-main {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-width: 0;
}
.chatgpt-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 20px;
  border-bottom: 1px solid #e5e7eb;
  background: #fff;
}
.chatgpt-header__left {
  display: flex;
  align-items: center;
  gap: 10px;
}
.chatgpt-header__left h2 {
  margin: 0;
  font-size: 18px;
}
.chatgpt-header__controls {
  display: flex;
  gap: 10px;
  align-items: center;
}
.chatgpt-scroll {
  flex: 1;
  overflow-y: auto;
  padding: 24px 20px;
}
.chatgpt-empty {
  max-width: 640px;
  margin: 40px auto;
  text-align: center;
}
.chatgpt-empty__icon {
  font-size: 48px;
  margin-bottom: 12px;
}
.chatgpt-suggestions {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  justify-content: center;
  margin-top: 20px;
}
.chatgpt-suggestion {
  border: 1px solid #d1d5db;
  background: #fff;
  border-radius: 999px;
  padding: 8px 14px;
  cursor: pointer;
  font-size: 13px;
}
.chatgpt-suggestion:hover {
  border-color: #10a37f;
  color: #10a37f;
}
.chatgpt-row {
  display: flex;
  gap: 12px;
  margin-bottom: 20px;
  max-width: 900px;
  margin-left: auto;
  margin-right: auto;
}
.chatgpt-row--user {
  flex-direction: row-reverse;
}
.chatgpt-avatar {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  background: #10a37f;
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 13px;
  flex-shrink: 0;
}
.chatgpt-avatar--user {
  background: #6366f1;
}
.chatgpt-bubble {
  flex: 1;
  min-width: 0;
  background: #fff;
  border: 1px solid #e5e7eb;
  border-radius: 12px;
  padding: 12px 16px;
}
.chatgpt-row--user .chatgpt-bubble {
  background: #eef2ff;
  border-color: #c7d2fe;
}
.chatgpt-markdown :deep(pre) {
  background: #1e293b;
  color: #f8fafc;
  padding: 12px;
  border-radius: 8px;
  overflow-x: auto;
}
.chatgpt-markdown :deep(code) {
  font-family: ui-monospace, monospace;
  font-size: 13px;
}
.chatgpt-meta {
  margin-top: 8px;
  font-size: 12px;
  color: #6b7280;
}
.chatgpt-thinking span {
  display: inline-block;
  width: 8px;
  height: 8px;
  margin-right: 4px;
  background: #9ca3af;
  border-radius: 50%;
  animation: blink 1.2s infinite;
}
.chatgpt-thinking span:nth-child(2) { animation-delay: 0.2s; }
.chatgpt-thinking span:nth-child(3) { animation-delay: 0.4s; }
@keyframes blink {
  0%, 80%, 100% { opacity: 0.3; }
  40% { opacity: 1; }
}
.chatgpt-footer {
  padding: 16px 20px;
  border-top: 1px solid #e5e7eb;
  background: #fff;
}
.chatgpt-input-wrap {
  max-width: 900px;
  margin: 0 auto;
  display: flex;
  gap: 10px;
  align-items: flex-end;
}
.chatgpt-input-wrap :deep(.el-textarea) {
  flex: 1;
}
</style>
