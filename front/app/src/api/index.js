const API_BASE = '/api/v1'
const TOKEN_KEY = 'ai-learning-system:token'
const SESSION_KEY = 'ai-learning-system:session'

async function request(path, options = {}) {
  const token = localStorage.getItem(TOKEN_KEY)
  const res = await fetch(`${API_BASE}${path}`, {
    ...options,
    headers: {
      'Content-Type': 'application/json',
      ...(token ? { Authorization: `Bearer ${token}` } : {}),
      ...options.headers
    }
  })
  const data = await res.json()
  if (data.code !== 0) {
    throw new Error(data.message || '请求失败')
  }
  return data.data
}

function authHeaders() {
  const token = localStorage.getItem(TOKEN_KEY)
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

export const authApi = {
  register: (body) => request('/auth/register', { method: 'POST', body: JSON.stringify(body) }),
  login: (username, password) => request('/auth/login', { method: 'POST', body: JSON.stringify({ username, password }) }),
  guest: () => request('/auth/guest', { method: 'POST' }),
  me: () => request('/auth/me'),
  permissions: () => request('/auth/permissions'),
  refreshPermissions: () => request('/auth/permissions/refresh', { method: 'POST' })
}

export const aiApi = {
  config: () => request('/ai/config'),
  chat: (question, history = []) =>
    request('/ai/chat', { method: 'POST', body: JSON.stringify({ question, history }) }),
  chatStream(question, history, handlers = {}) {
    const controller = new AbortController()
    const run = async () => {
      const res = await fetch(`${API_BASE}/ai/chat/stream`, {
        method: 'POST',
        headers: authHeaders(),
        body: JSON.stringify({ question, history }),
        signal: controller.signal
      })
      await consumeSSE(res, handlers)
    }
    run().catch((err) => {
      if (err.name !== 'AbortError' && handlers.onError) handlers.onError(err)
    })
    return () => controller.abort()
  },
  careerInterview(message, history, handlers = {}) {
    const controller = new AbortController()
    const run = async () => {
      const res = await fetch(`${API_BASE}/ai/career/interview`, {
        method: 'POST',
        headers: authHeaders(),
        body: JSON.stringify({ message, history }),
        signal: controller.signal
      })
      await consumeSSE(res, handlers)
    }
    run().catch((err) => {
      if (err.name !== 'AbortError' && handlers.onError) handlers.onError(err)
    })
    return () => controller.abort()
  },
  careerRecommend: (body) => request('/ai/career/recommend', { method: 'POST', body: JSON.stringify(body) }),
  goalDecompose: (body) => request('/ai/goal/decompose', { method: 'POST', body: JSON.stringify(body) }),
  learningSuggest: (body) => request('/ai/learning/suggest', { method: 'POST', body: JSON.stringify(body) })
}

export function saveSession(token, user) {
  localStorage.setItem(TOKEN_KEY, token)
  localStorage.setItem(SESSION_KEY, JSON.stringify({ ...user, isGuest: user.role === 'guest' }))
}

export function clearSession() {
  localStorage.removeItem(TOKEN_KEY)
  localStorage.removeItem(SESSION_KEY)
}

export function loadSession() {
  try {
    const raw = localStorage.getItem(SESSION_KEY)
    return raw ? JSON.parse(raw) : null
  } catch {
    return null
  }
}
