package mail

import "errors"

var (
	ErrNotFound     = errors.New("mail: message not found")
	ErrNoRecipients = errors.New("mail: email has no recipients")
)
