package mail

import "time"

type Message struct {
	ID          string       `json:"id"`
	From        string       `json:"from"`
	To          []string     `json:"to"`
	Subject     string       `json:"subject"`
	TextBody    string       `json:"text_body"`
	HTMLBody    string       `json:"html_body"`
	Raw         string       `json:"raw"`
	Attachments []Attachment `json:"attachments"`
	CreatedAt   time.Time    `json:"created_at"`
}

func (m *Message) HasRecipients() bool {
	return len(m.To) >= 1
}

type Attachment struct {
	Filename    string `json:"filename"`
	ContentType string `json:"content_type"`
	Content     []byte `json:"content"`
}
