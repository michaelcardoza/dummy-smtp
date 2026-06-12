import hljs from 'highlight.js/lib/core'
import python from 'highlight.js/lib/languages/python'
import go from 'highlight.js/lib/languages/go'
import javascript from 'highlight.js/lib/languages/javascript'
import php from 'highlight.js/lib/languages/php'
import ruby from 'highlight.js/lib/languages/ruby'
import bash from 'highlight.js/lib/languages/bash'
import xml from 'highlight.js/lib/languages/xml'
import json from 'highlight.js/lib/languages/json'

hljs.registerLanguage('python', python)
hljs.registerLanguage('go', go)
hljs.registerLanguage('javascript', javascript)
hljs.registerLanguage('php', php)
hljs.registerLanguage('ruby', ruby)
hljs.registerLanguage('bash', bash)
hljs.registerLanguage('xml', xml)
hljs.registerLanguage('json', json)

/** Syntax-highlight code to HTML using the given highlight.js grammar. */
export function highlight(code: string, language: string): string {
  return hljs.highlight(code, { language }).value
}
