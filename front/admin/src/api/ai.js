import api from './index.js'

const API_BASE = '/api/v1'

function authHeaders() {
  const token = localStorage.getItem('admin-token')
  return {
    'Content-Type': 'application/json',
    ...(token ? { Authorization: `Bearer ${token}` } : {})
  }
}

async function consumeSSE(res, { onToken, onDone, onError }) {
  if (!res.ok) {
    const text = await res.text()
    throw new Error(text || `请求失败 (${res.status})`)
  }
  const reader = res.body.getReader()
  const decoder = new TextDecoder()
  let buffer = ''

  while (true) {
    const { done, value } = await reader.read()
    if (done) break
    buffer += decoder.decode(value, { stream: true })
    const parts = buffer.split('\n\n')
    buffer = parts.pop() || ''
    for (const part of parts) {
      const line = part.split('\n').find((l) => l.startsWith('data: '))
      if (!line) continue
      let payload
      try {
        payload = JSON.parse(line.slice(6))
      } catch {
        continue
      }
      if (payload.type === 'token' && onToken) onToken(payload.content || '')
      if (payload.type === 'error') throw new Error(payload.message || 'AI 服务错误')
      if (payload.type === 'done' && onDone) onDone(payload)
    }
  }
}

export const aiApi = {
  config: () => api.get('/ai/config'),
  chatStream(question, history, options = {}, handlers = {}) {
    const controller = new AbortController()
    const run = async () => {
      const res = await fetch(`${API_BASE}/ai/chat/stream`, {
        method: 'POST',
        headers: authHeaders(),
        body: JSON.stringify({
          question,
          history,
          virtualModel: options.virtualModel || '',
          mode: options.mode || 'chat'
        }),
        signal: controller.signal
      })
      await consumeSSE(res, handlers)
    }
    run().catch((err) => {
      if (err.name !== 'AbortError' && handlers.onError) handlers.onError(err)
    })
    return () => controller.abort()
  }
}
