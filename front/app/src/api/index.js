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

export const authApi = {
  register: (body) => request('/auth/register', { method: 'POST', body: JSON.stringify(body) }),
  login: (username, password) => request('/auth/login', { method: 'POST', body: JSON.stringify({ username, password }) }),
  guest: () => request('/auth/guest', { method: 'POST' }),
  me: () => request('/auth/me'),
  permissions: () => request('/auth/permissions')
}

export const aiApi = {
  chat: (question) => request('/ai/chat', { method: 'POST', body: JSON.stringify({ question }) }),
  config: () => request('/ai/config')
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
