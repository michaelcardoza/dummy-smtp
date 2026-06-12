package web

import (
	"sync"

	"github.com/michaelcardoza/dummy-smtp/internal/core/mail"
)

type Broker struct {
	mu   sync.Mutex
	subs map[chan *mail.Message]struct{}
}

func NewBroker() *Broker {
	return &Broker{
		subs: map[chan *mail.Message]struct{}{},
	}
}

func (b *Broker) Publish(message *mail.Message) {
	b.mu.Lock()
	defer b.mu.Unlock()
	for ch := range b.subs {
		select {
		case ch <- message:
		default:
		}
	}
}

func (b *Broker) Subscribe() chan *mail.Message {
	ch := make(chan *mail.Message, 8)
	b.mu.Lock()
	b.subs[ch] = struct{}{}
	b.mu.Unlock()
	return ch
}

func (b *Broker) Unsubscribe(ch chan *mail.Message) {
	b.mu.Lock()
	delete(b.subs, ch)
	b.mu.Unlock()
	close(ch)
}
