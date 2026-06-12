import type { Message } from '../types'
import { fileStamp } from './time'

/** Lowercase, hyphenated, filesystem-safe slug. */
export function slugify(text: string, maxLength = 50): string {
  return text
    .toLowerCase()
    .replace(/[^a-z0-9]+/g, '-')
    .replace(/^-+|-+$/g, '')
    .slice(0, maxLength)
}

/** Download base name for a message: "2026-06-12_1430-subject-slug". */
export function downloadBaseName(message: Message): string {
  return `${fileStamp(message.created_at)}-${slugify(message.subject) || 'no-subject'}`
}
