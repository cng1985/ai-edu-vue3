import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { customerApi } from '../api'
import { useWebSocket } from '../composables/useWebSocket'

export const useCustomerStore = defineStore('customer', () => {
  const tickets = ref([])
  const activeTicketId = ref(null)
  const messages = ref([])
  const loading = ref(false)
  const sending = ref(false)

  const ws = useWebSocket()
  let unsubscribers = []

  const activeTicket = computed(() => tickets.value.find((t) => t.id === activeTicketId.value))

  function setupWS() {
    ws.connect()
    unsubscribers.push(
      ws.on('support.message', (msg) => {
        if (msg.ticketId === activeTicketId.value) {
          if (!messages.value.find((m) => m.id === msg.id)) {
            messages.value.push(msg)
          }
        }
        const ticket = tickets.value.find((t) => t.id === msg.ticketId)
        if (ticket) {
          ticket.lastMessage = msg.content
          ticket.lastMessageAt = msg.createdAt
        }
      }),
      ws.on('support.ticket.update', (ticket) => {
        const idx = tickets.value.findIndex((t) => t.id === ticket.id)
        if (idx >= 0) tickets.value[idx] = { ...tickets.value[idx], ...ticket }
        if (activeTicketId.value === ticket.id && ticket.status === 'closed') {
          // refresh messages view
        }
      })
    )
  }

  function teardownWS() {
    unsubscribers.forEach((fn) => fn())
    unsubscribers = []
    ws.disconnect()
  }

  async function loadTickets() {
    loading.value = true
    try {
      const res = await customerApi.listTickets()
      tickets.value = res.list || []
    } finally {
      loading.value = false
    }
  }

  async function selectTicket(id) {
    activeTicketId.value = id
    loading.value = true
    try {
      const res = await customerApi.listMessages(id)
      messages.value = res.list || []
      ws.send('support.subscribe', { ticketId: id })
    } finally {
      loading.value = false
    }
  }

  async function createTicket(subject, content) {
    sending.value = true
    try {
      const res = await customerApi.createTicket({ subject, content })
      tickets.value.unshift(res.ticket)
      activeTicketId.value = res.ticket.id
      messages.value = [res.message]
      ws.send('support.subscribe', { ticketId: res.ticket.id })
      return res.ticket
    } finally {
      sending.value = false
    }
  }

  async function sendMessage(content) {
    if (!activeTicketId.value || !content.trim()) return
    sending.value = true
    const sent = ws.send('support.send', { ticketId: activeTicketId.value, content: content.trim() })
    if (!sent) {
      try {
        const msg = await customerApi.sendMessage(activeTicketId.value, content.trim())
        messages.value.push(msg)
      } finally {
        sending.value = false
      }
      return
    }
    sending.value = false
  }

  function reset() {
    activeTicketId.value = null
    messages.value = []
    tickets.value = []
  }

  return {
    tickets,
    activeTicketId,
    activeTicket,
    messages,
    loading,
    sending,
    connected: ws.connected,
    setupWS,
    teardownWS,
    loadTickets,
    selectTicket,
    createTicket,
    sendMessage,
    reset
  }
})
