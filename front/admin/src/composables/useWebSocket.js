import { ref, onUnmounted } from 'vue'

const ADMIN_TOKEN_KEY = 'admin-token'

function adminWsURL(path) {
  const proto = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
  const host = window.location.host
  const token = localStorage.getItem(ADMIN_TOKEN_KEY)
  return `${proto}//${host}${path}?token=${encodeURIComponent(token || '')}`
}

export function useAdminWebSocket(path = '/api/v1/ws/support') {
  const connected = ref(false)
  const listeners = new Map()
  let ws = null
  let reconnectTimer = null
  let reconnectDelay = 1000
  let intentionalClose = false

  function on(event, handler) {
    if (!listeners.has(event)) listeners.set(event, new Set())
    listeners.get(event).add(handler)
    return () => listeners.get(event)?.delete(handler)
  }

  function emit(event, payload) {
    listeners.get(event)?.forEach((fn) => fn(payload))
  }

  function connect() {
    if (ws && (ws.readyState === WebSocket.OPEN || ws.readyState === WebSocket.CONNECTING)) return
    intentionalClose = false
    ws = new WebSocket(adminWsURL(path))

    ws.onopen = () => {
      connected.value = true
      reconnectDelay = 1000
    }

    ws.onmessage = (e) => {
      try {
        const { event, payload } = JSON.parse(e.data)
        const data = typeof payload === 'string' ? JSON.parse(payload) : payload
        emit(event, data)
      } catch { /* ignore */ }
    }

    ws.onclose = () => {
      connected.value = false
      if (!intentionalClose) {
        reconnectTimer = setTimeout(() => {
          reconnectTimer = null
          reconnectDelay = Math.min(reconnectDelay * 2, 30000)
          connect()
        }, reconnectDelay)
      }
    }
  }

  function send(event, payload) {
    if (!ws || ws.readyState !== WebSocket.OPEN) return false
    ws.send(JSON.stringify({ event, payload }))
    return true
  }

  function disconnect() {
    intentionalClose = true
    if (reconnectTimer) clearTimeout(reconnectTimer)
    if (ws) { ws.close(); ws = null }
    connected.value = false
  }

  onUnmounted(disconnect)

  return { connected, connect, disconnect, send, on }
}
