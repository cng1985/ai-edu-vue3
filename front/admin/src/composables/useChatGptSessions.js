import { ref, computed } from 'vue'

const STORAGE_KEY = 'admin-chatgpt-sessions'

function loadSessions() {
  try {
    const raw = localStorage.getItem(STORAGE_KEY)
    if (raw) return JSON.parse(raw)
  } catch { /* ignore */ }
  return []
}

function saveSessions(list) {
  localStorage.setItem(STORAGE_KEY, JSON.stringify(list))
}

export function useChatGptSessions() {
  const sessions = ref(loadSessions())
  const activeId = ref(sessions.value[0]?.id || '')

  const activeSession = computed(() =>
    sessions.value.find((s) => s.id === activeId.value) || null
  )

  function ensureActive() {
    if (!sessions.value.length) {
      createSession()
    } else if (!activeId.value || !activeSession.value) {
      activeId.value = sessions.value[0].id
    }
  }

  function createSession() {
    const session = {
      id: `s_${Date.now()}`,
      title: '新对话',
      messages: [],
      virtualModel: '',
      mode: 'chat',
      updatedAt: Date.now()
    }
    sessions.value.unshift(session)
    activeId.value = session.id
    persist()
    return session
  }

  function selectSession(id) {
    activeId.value = id
  }

  function deleteSession(id) {
    const idx = sessions.value.findIndex((s) => s.id === id)
    if (idx === -1) return
    sessions.value.splice(idx, 1)
    if (activeId.value === id) {
      activeId.value = sessions.value[0]?.id || ''
      if (!sessions.value.length) createSession()
    }
    persist()
  }

  function updateSession(id, patch) {
    const s = sessions.value.find((x) => x.id === id)
    if (!s) return
    Object.assign(s, patch)
    s.updatedAt = Date.now()
    persist()
  }

  function setMessages(id, messages) {
    updateSession(id, { messages })
    const firstUser = messages.find((m) => m.role === 'user')
    if (firstUser?.text) {
      const title = firstUser.text.slice(0, 32) + (firstUser.text.length > 32 ? '…' : '')
      updateSession(id, { title })
    }
  }

  function persist() {
    saveSessions(sessions.value)
  }

  ensureActive()

  return {
    sessions,
    activeId,
    activeSession,
    createSession,
    selectSession,
    deleteSession,
    updateSession,
    setMessages,
    ensureActive
  }
}
