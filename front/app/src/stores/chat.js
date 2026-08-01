import { defineStore } from 'pinia'
import { ask } from '../utils/aiEngine'

const STORAGE_KEY = 'ai-learning-system:chat'

function loadMessages() {
  try {
    const raw = localStorage.getItem(STORAGE_KEY)
    if (raw) return JSON.parse(raw)
  } catch (e) { /* 忽略损坏的本地数据 */ }
  return []
}

let cancelStream = null

export const useChatStore = defineStore('chat', {
  state: () => ({
    // { id, role: 'user' | 'assistant', text, sources, streaming }
    messages: loadMessages(),
    generating: false
  }),

  actions: {
    persist() {
      // 只持久化已完成的消息
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

      cancelStream = ask(text, {
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
