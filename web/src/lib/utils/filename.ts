import type { Message } from '../types'
import { fileStamp } from './time'

export function slugify(text: string, maxLength = 50): string {
  return text
    .toLowerCase()
    .replace(/[^a-z0-9]+/g, '-')
    .replace(/^-+|-+$/g, '')
    .slice(0, maxLength)
}

export function downloadBaseName(message: Message): string {
  return `${fileStamp(message.createdAt)}-${slugify(message.subject) || 'no-subject'}`
}
