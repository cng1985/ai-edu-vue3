import { defineStore } from 'pinia'
import { aiApi } from '../api'
import { ask as localAsk } from '../utils/aiEngine'

const STORAGE_KEY = 'ai-learning-system:chat'

function loadMessages() {
  try {
    const raw = localStorage.getItem(STORAGE_KEY)
    if (raw) return JSON.parse(raw)
  } catch (e) { /* 忽略损坏的本地数据 */ }
  return []
}

function buildHistory(messages) {
  return messages
    .filter((m) => !m.streaming && m.text)
    .slice(-16)
    .map((m) => ({ role: m.role, content: m.text }))
}

let cancelStream = null

export const useChatStore = defineStore('chat', {
  state: () => ({
    messages: loadMessages(),
    generating: false,
    aiConfig: null,
    useBackend: true
  }),

  actions: {
    async loadConfig() {
      try {
        this.aiConfig = await aiApi.config()
      } catch {
        this.aiConfig = { enabled: false, model: 'local', baseUrl: '' }
      }
    },

    persist() {
      const done = this.messages.filter((m) => !m.streaming)
      localStorage.setItem(STORAGE_KEY, JSON.stringify(done))
    },

    send(question) {
      const text = question.trim()
      if (!text || this.generating) return

      this.messages.push({
        id: Date.now() + ':user',
        role: 'user',
        text,
        sources: [],
        streaming: false
      })

      const reply = {
        id: Date.now() + ':assistant',
        role: 'assistant',
        text: '',
        sources: [],
        streaming: true
      }
      this.messages.push(reply)
      this.generating = true
      this.persist()

      const history = buildHistory(this.messages.slice(0, -2))
      const handlers = {
        onToken: (chunk) => {
          reply.text += chunk
        },
        onDone: (result) => {
          reply.text = result.text || reply.text
          reply.sources = result.sources || []
          reply.streaming = false
          this.generating = false
          cancelStream = null
          this.persist()
        },
        onError: () => {
          this.fallbackLocal(text, reply)
        }
      }

      if (this.useBackend) {
        cancelStream = aiApi.chatStream(text, history, handlers)
      } else {
        this.fallbackLocal(text, reply)
      }
    },

    fallbackLocal(question, reply) {
      cancelStream = localAsk(question, {
        onToken: (chunk) => {
          reply.text += chunk
        },
        onDone: (result) => {
          reply.text = result.text
          reply.sources = result.sources
          reply.streaming = false
          this.generating = false
          cancelStream = null
          this.persist()
        }
      })
    },

    stop() {
      if (cancelStream) {
        cancelStream()
        cancelStream = null
      }
      const last = this.messages[this.messages.length - 1]
      if (last && last.streaming) {
        last.streaming = false
        last.text += '\n\n*（已停止生成）*'
      }
      this.generating = false
      this.persist()
    },

    clear() {
      this.stop()
      this.messages = []
      this.persist()
    }
  }
})
