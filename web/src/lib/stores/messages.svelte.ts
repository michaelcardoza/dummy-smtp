import type { Message } from '../types'
import { listMessages, deleteMessage, deleteAllMessages, sendTestMessage } from '../api'
import { toasts } from './toasts.svelte'
import { EVENTS_URL } from '../config'

class MessagesStore {
  items = $state<Message[]>([])
  selectedId = $state<string | null>(null)
  searchQuery = $state('')
  isConnected = $state(false)
  liveEnabled = $state(true)
  private eventSource: EventSource | null = null

  readonly filtered = $derived.by(() => {
    const query = this.searchQuery.trim().toLowerCase()
    const matches = this.items.filter((message) => {
      if (!query) return true
      return (
        message.subject.toLowerCase().includes(query) ||
        message.from.toLowerCase().includes(query) ||
        message.to.join(', ').toLowerCase().includes(query)
      )
    })
    return matches.sort(
      (a, b) => new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime(),
    )
  })

  readonly selected = $derived(this.items.find((m) => m.id === this.selectedId) ?? null)
  readonly filteredCount = $derived(this.filtered.length)
  readonly totalCount = $derived(this.items.length)

  load = async () => {
    try {
      this.items = await listMessages()
    } catch {
      toasts.push('Failed to load messages', true)
    }
  }

  select = (id: string) => {
    this.selectedId = id
  }

  receive = (message: Message) => {
    if (this.items.some((m) => m.id === message.id)) return
    this.items = [message, ...this.items]
    toasts.push(`New mail · ${message.subject || '(no subject)'}`)
  }

  remove = async (id: string) => {
    try {
      await deleteMessage(id)
      this.items = this.items.filter((m) => m.id !== id)
      if (this.selectedId === id) this.selectedId = null
    } catch {
      toasts.push('Failed to delete message', true)
    }
  }

  clear = async () => {
    try {
      await deleteAllMessages()
      this.items = []
      this.selectedId = null
    } catch {
      toasts.push('Failed to clear messages', true)
    }
  }

  sendTest = async () => {
    try {
      await sendTestMessage()
      toasts.push('Test message sent')
    } catch {
      toasts.push('Failed to send test', true)
    }
  }

  connectLive = () => {
    if (this.eventSource) return
    const source = new EventSource(EVENTS_URL)
    source.onopen = () => (this.isConnected = true)
    source.onerror = () => (this.isConnected = false)
    source.addEventListener('message', (e) => this.receive(JSON.parse(e.data) as Message))
    this.eventSource = source
  }

  disconnectLive = () => {
    this.eventSource?.close()
    this.eventSource = null
    this.isConnected = false
  }

  toggleLive = () => {
    this.liveEnabled = !this.liveEnabled
    if (this.liveEnabled) this.connectLive()
    else this.disconnectLive()
  }
}

export const messages = new MessagesStore()
