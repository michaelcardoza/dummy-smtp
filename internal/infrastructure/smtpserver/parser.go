package smtpserver

import (
	"io"
	"mime"
	"mime/multipart"
	"net/mail"
	"strings"

	coremail "github.com/michaelcardoza/dummy-smtp/internal/core/mail"
)

type parsedMessage struct {
	Subject     string
	TextBody    string
	HTMLBody    string
	Attachments []coremail.Attachment
}

func parse(raw string) parsedMessage {
	var p parsedMessage

	msg, err := mail.ReadMessage(strings.NewReader(raw))
	if err != nil {
		return p
	}

	p.Subject = decodeHeader(msg.Header.Get("Subject"))

	mediaType, params, err := mime.ParseMediaType(msg.Header.Get("Content-Type"))
	if err != nil {
		body, _ := io.ReadAll(msg.Body)
		p.TextBody = string(body)
		return p
	}

	if strings.HasPrefix(mediaType, "multipart/") {
		parseMultipart(msg.Body, params["boundary"], &p)
		return p
	}

	body, _ := io.ReadAll(msg.Body)
	if mediaType == "text/html" {
		p.HTMLBody = string(body)
	} else {
		p.TextBody = string(body)
	}

	return p
}

func parseMultipart(body io.Reader, boundary string, p *parsedMessage) {
	if boundary == "" {
		return
	}
	mr := multipart.NewReader(body, boundary)
	for {
		part, err := mr.NextPart()
		if err != nil {
			return
		}
		mediaType, params, _ := mime.ParseMediaType(part.Header.Get("Content-Type"))

		if strings.HasPrefix(mediaType, "multipart/") {
			parseMultipart(part, params["boundary"], p)
			continue
		}

		content, _ := io.ReadAll(part)

		if filename := part.FileName(); filename != "" {
			p.Attachments = append(p.Attachments, coremail.Attachment{
				Filename:    filename,
				ContentType: mediaType,
				Size:        len(content),
			})
			continue
		}

		switch mediaType {
		case "text/html":
			p.HTMLBody = string(content)
		case "text/plain":
			p.TextBody = string(content)
		}
	}
}

func decodeHeader(s string) string {
	dec := new(mime.WordDecoder)
	out, err := dec.DecodeHeader(s)
	if err != nil {
		return s
	}
	return out
}
