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
  let data
  try {
    data = await res.json()
  } catch {
    throw new Error(`请求失败 (${res.status})`)
  }
  if (data.code !== 0) {
    throw new Error(data.message || '请求失败')
  }
  return data.data
}

const get = (path, params) => {
  const qs = params ? new URLSearchParams(Object.entries(params).filter(([, v]) => v !== undefined && v !== '')) : null
  return request(qs && String(qs) ? `${path}?${qs}` : path)
}
const post = (path, body) => request(path, { method: 'POST', body: JSON.stringify(body || {}) })
const put = (path, body) => request(path, { method: 'PUT', body: JSON.stringify(body || {}) })
const del = (path) => request(path, { method: 'DELETE' })

function authHeaders() {
  const token = localStorage.getItem(TOKEN_KEY)
  return {
    'Content-Type': 'application/json',
    ...(token ? { Authorization: `Bearer ${token}` } : {})
  }
}

async function consumeSSE(res, handlers) {
  if (!res.ok) {
    const text = await res.text()
    let message = text
    try {
      message = JSON.parse(text).message || text
    } catch { /* 非 JSON 错误体 */ }
    throw new Error(message || `请求失败 (${res.status})`)
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
      if (payload.type === 'error') throw new Error(payload.message || 'AI 服务错误')
      handlers[payload.type]?.(payload)
    }
  }
}

/**
 * 发起 SSE 流式请求。handlers 以事件类型为键：token / stage / done / error。
 * 返回取消函数。
 */
function stream(path, body, handlers = {}) {
  const controller = new AbortController()
  fetch(`${API_BASE}${path}`, {
    method: 'POST',
    headers: authHeaders(),
    body: JSON.stringify(body || {}),
    signal: controller.signal
  })
    .then((res) => consumeSSE(res, handlers))
    .catch((err) => {
      if (err.name !== 'AbortError') handlers.error?.(err)
    })
  return () => controller.abort()
}

export const authApi = {
  register: (body) => post('/auth/register', body),
  login: (username, password) => post('/auth/login', { username, password }),
  guest: () => post('/auth/guest'),
  me: () => get('/auth/me'),
  permissions: () => get('/auth/permissions'),
  refreshPermissions: () => post('/auth/permissions/refresh')
}

export const aiApi = {
  config: () => get('/ai/config')
}

/** 职业体系、知识图谱、项目、任务市场、社区与知识资产 */
export const ecoApi = {
  careers: () => get('/eco/careers'),
  role: (id) => get(`/eco/roles/${id}`),
  skills: () => get('/eco/skills'),
  graph: (skillId) => get('/eco/knowledge/graph', { skillId }),
  knowledge: (id) => get(`/eco/knowledge/${id}`),
  projects: () => get('/eco/projects'),
  project: (id) => get(`/eco/projects/${id}`),
  submit: (projectId, taskId, body) => post(`/eco/projects/${projectId}/tasks/${taskId}/submissions`, body),
  opportunities: (params) => get('/eco/opportunities', params),
  opportunity: (id) => get(`/eco/opportunities/${id}`),
  apply: (id, message) => post(`/eco/opportunities/${id}/apply`, { message }),
  posts: (params) => get('/eco/community/posts', params),
  post: (id) => get(`/eco/community/posts/${id}`),
  createPost: (body) => post('/eco/community/posts', body),
  answer: (id, content) => post(`/eco/community/posts/${id}/answers`, { content }),
  likePost: (id) => post(`/eco/community/posts/${id}/like`),
  likeAnswer: (id) => post(`/eco/community/answers/${id}/like`),
  acceptAnswer: (id) => post(`/eco/community/answers/${id}/accept`),
  summarize: (id) => post(`/eco/community/posts/${id}/summarize`),
  resources: (params) => get('/eco/resources', params),
  resource: (id) => get(`/eco/resources/${id}`),
  createResource: (body) => post('/eco/resources', body)
}

/** 个人成长数据 */
export const meApi = {
  overview: () => get('/me/overview'),
  goal: () => get('/me/goal'),
  setGoal: (body) => put('/me/goal', body),
  gap: (roleId) => get('/me/gap', { roleId }),
  skills: () => get('/me/skills'),
  knowledgeStates: () => get('/me/knowledge-states'),
  recommendations: (limit) => get('/me/recommendations', { limit }),
  reviews: () => get('/me/reviews'),
  profile: () => get('/me/profile'),
  submissions: () => get('/me/submissions'),
  applications: () => get('/me/applications'),
  recordEvent: (body) => post('/me/events', body),
  practice: (knowledgeId, answers) => post(`/me/knowledge/${knowledgeId}/practice`, { answers }),
  completeChapter: (body) => post('/me/chapters/complete', body)
}

/** 企业：发布任务、处理申请、推荐人才 */
export const enterpriseApi = {
  published: () => get('/enterprise/opportunities'),
  publish: (body) => post('/enterprise/opportunities', body),
  update: (id, body) => put(`/enterprise/opportunities/${id}`, body),
  remove: (id) => del(`/enterprise/opportunities/${id}`),
  candidates: (id) => get(`/enterprise/opportunities/${id}/candidates`),
  decide: (applicationId, body) => post(`/enterprise/applications/${applicationId}/decision`, body)
}

export const talentApi = {
  list: (keyword) => get('/talents', { keyword }),
  profile: (id) => get(`/talents/${id}`)
}

/** AI Learning Kernel 与 Agent */
export const kernelApi = {
  stages: () => get('/ai/kernel/stages'),
  runs: () => get('/ai/kernel/runs'),
  runStream: (question, handlers) => stream('/ai/kernel/run/stream', { question }, handlers),
  agents: () => get('/ai/agents'),
  agentContext: (code, message) => get(`/ai/agents/${code}/context`, { message }),
  agentChat: (code, message, history, handlers) => stream(`/ai/agents/${code}/chat/stream`, { message, history }, handlers)
}

export const customerApi = {
  listTickets: (params) => get('/app/support/tickets', params),
  createTicket: (body) => post('/app/support/tickets', body),
  getTicket: (id) => get(`/app/support/tickets/${id}`),
  listMessages: (id, params) => get(`/app/support/tickets/${id}/messages`, params),
  sendMessage: (id, content) => post(`/app/support/tickets/${id}/messages`, { content })
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
