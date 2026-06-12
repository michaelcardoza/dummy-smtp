package mail

type Publisher interface {
	Publish(message *Message)
}
