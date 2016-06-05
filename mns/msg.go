package mns

type Message struct {
	BotID string
	Type  string
	Body  interface{}
}
