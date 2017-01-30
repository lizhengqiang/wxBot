package domain

type Message struct {
	BotID string
	Type  string
	Body  interface{}
}

type MessageHandler func(*Message) error

type MessageQueue interface {
	Send(*Message) error
	RegisterHandler(MessageHandler) (error)
}
