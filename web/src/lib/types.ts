export interface Attachment {
  filename: string
  content_type: string
}

export interface Message {
  id: string
  from: string
  to: string[]
  subject: string
  textBody: string
  htmlBody: string
  raw: string
  attachments: Attachment[]
  createdAt: string
}
