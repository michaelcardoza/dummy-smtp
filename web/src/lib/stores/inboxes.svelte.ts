import type { Inbox } from '../types'

const KEY = 'dummy-smtp:inboxes'

const seed: Inbox[] = [
  { id: '1', name: 'Project A', domain: 'projecta.test', description: 'Staging mail for Project A' },
  { id: '2', name: 'Project B', domain: 'projectb.test', description: '' },
]

class InboxStore {
  items = $state<Inbox[]>(read())

  create = (data: Omit<Inbox, 'id'>) => {
    const inbox: Inbox = { id: crypto.randomUUID(), ...data }
    this.items = [...this.items, inbox]
    this.persist()
    return inbox
  }

  update = (id: string, data: Omit<Inbox, 'id'>) => {
    this.items = this.items.map((i) => (i.id === id ? { ...i, ...data } : i))
    this.persist()
  }

  remove = (id: string) => {
    this.items = this.items.filter((i) => i.id !== id)
    this.persist()
  }

  private persist() {
    localStorage.setItem(KEY, JSON.stringify(this.items))
  }
}

function read(): Inbox[] {
  try {
    const raw = localStorage.getItem(KEY)
    return raw ? (JSON.parse(raw) as Inbox[]) : seed
  } catch {
    return seed
  }
}

export const inboxes = new InboxStore()
