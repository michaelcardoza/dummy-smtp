export interface Attachment {
  filename: string
  content_type: string
}

export interface Message {
  id: string
  from: string
  to: string[]
  subject: string
  text_body: string
  html_body: string
  raw: string
  attachments: Attachment[]
  created_at: string
}

export interface Inbox {
  id: string
  name: string
  domain: string
  description: string
}
