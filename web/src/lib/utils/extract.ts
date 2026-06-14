import type { Message } from '../types'

function stripTags(html: string): string {
  return html.replace(/<[^>]+>/g, ' ')
}

export function extractLinks(m: Message): string[] {
  const text = `${m.textBody}\n${m.htmlBody}`
  const urls = text.match(/https?:\/\/[^\s"'<>)]+/g) ?? []
  return [...new Set(urls)]
}

export function extractCodes(m: Message): string[] {
  const text = m.textBody || stripTags(m.htmlBody)
  const codes = text.match(/\b\d{4,8}\b/g) ?? []
  return [...new Set(codes)]
}
