import type { Inbox, Message } from '../types'
import { listMessages, deleteMessage, deleteAllMessages, sendTestMessage } from '../api'
import { toasts } from './toasts.svelte'
import { EVENTS_URL } from '../config'

/** Domain of a message's first recipient — used to route messages to inboxes. */
export function recipientDomain(message: Message): string {
  const recipient = message.to[0] ?? ''
  const at = recipient.indexOf('@')
  return at === -1 ? '' : recipient.slice(at + 1)
}

class MessagesStore {
  items = $state<Message[]>([])
  selectedId = $state<string | null>(null)
  activeInbox = $state<Inbox | null>(null)
  searchQuery = $state('')
  isConnected = $state(false)
  liveEnabled = $state(true)
  private eventSource: EventSource | null = null

  /** Messages after applying the active inbox and search filters. */
  readonly filtered = $derived.by(() => {
    const query = this.searchQuery.trim().toLowerCase()
    return this.items.filter((message) => {
      if (this.activeInbox && recipientDomain(message) !== this.activeInbox.domain) return false
      if (!query) return true
      return (
        message.subject.toLowerCase().includes(query) ||
        message.from.toLowerCase().includes(query) ||
        message.to.join(', ').toLowerCase().includes(query)
      )
    })
  })

  readonly selected = $derived(this.items.find((m) => m.id === this.selectedId) ?? null)
  readonly filteredCount = $derived(this.filtered.length)
  readonly totalCount = $derived(this.items.length)

  countForDomain = (domain: string) =>
    this.items.filter((m) => recipientDomain(m) === domain).length

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

  selectInbox = (inbox: Inbox | null) => {
    this.activeInbox = inbox
    this.selectedId = null
  }

  /** Prepend a message pushed over the live stream (ignoring duplicates). */
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
