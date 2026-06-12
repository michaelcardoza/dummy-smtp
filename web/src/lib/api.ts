import type { Message } from './types'
import { MESSAGES_API } from './config'

export async function listMessages(): Promise<Message[]> {
  const res = await fetch(MESSAGES_API)
  if (!res.ok) throw new Error(`list failed: ${res.status}`)
  return res.json()
}

export async function getMessage(id: string): Promise<Message> {
  const res = await fetch(`${MESSAGES_API}/${id}`)
  if (!res.ok) throw new Error(`get failed: ${res.status}`)
  return res.json()
}

export async function deleteMessage(id: string): Promise<void> {
  const res = await fetch(`${MESSAGES_API}/${id}`, { method: 'DELETE' })
  if (!res.ok) throw new Error(`delete failed: ${res.status}`)
}

export async function deleteAllMessages(): Promise<void> {
  const res = await fetch(MESSAGES_API, { method: 'DELETE' })
  if (!res.ok) throw new Error(`clear failed: ${res.status}`)
}

export async function sendTestMessage(): Promise<void> {
  const res = await fetch(`${MESSAGES_API}/test`, { method: 'POST' })
  if (!res.ok) throw new Error(`send test failed: ${res.status}`)
}
