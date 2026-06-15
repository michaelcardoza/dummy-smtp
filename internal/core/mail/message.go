package mail

import "time"

type Message struct {
	ID          string       `json:"id" bson:"id"`
	From        string       `json:"from" bson:"from"`
	To          []string     `json:"to" bson:"to"`
	Subject     string       `json:"subject" bson:"subject"`
	TextBody    string       `json:"textBody" bson:"textBody"`
	HTMLBody    string       `json:"htmlBody" bson:"htmlBody"`
	Raw         string       `json:"raw" bson:"raw"`
	Attachments []Attachment `json:"attachments" bson:"attachments"`
	CreatedAt   time.Time    `json:"createdAt" bson:"createdAt"`
}

func (m *Message) HasRecipients() bool {
	return len(m.To) >= 1
}

type Attachment struct {
	Filename    string `json:"filename" bson:"filename"`
	ContentType string `json:"contentType" bson:"contentType"`
	Size        int    `json:"size" bson:"size"`
}
