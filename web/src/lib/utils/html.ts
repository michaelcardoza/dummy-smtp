const VOID_ELEMENTS = /^<(area|base|br|col|embed|hr|img|input|link|meta|param|source|track|wbr)\b/i

/** Pretty-print HTML one tag per line with indentation (for source view). */
export function formatHtml(html: string): string {
  const lines = html
    .replace(/\r/g, '')
    .replace(/>\s+</g, '><')
    .replace(/></g, '>\n<')
    .split('\n')

  let depth = 0
  const formatted: string[] = []
  for (const rawLine of lines) {
    const line = rawLine.trim()
    if (!line) continue

    const isClosingTag = /^<\//.test(line)
    const isSelfContained = /^<[^/!?][^>]*>.*<\/[^>]+>$/.test(line)
    const isOpeningTag =
      /^<[^/!?]/.test(line) && !/\/>$/.test(line) && !VOID_ELEMENTS.test(line) && !isSelfContained

    if (isClosingTag) depth = Math.max(0, depth - 1)
    formatted.push('  '.repeat(depth) + line)
    if (isOpeningTag) depth++
  }
  return formatted.join('\n')
}
