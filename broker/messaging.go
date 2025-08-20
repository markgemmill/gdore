package broker

import "fmt"

type Message struct {
	Message string
	Done    bool
	Err     error
	File    string
}

type MessageBroker struct {
	Channel chan Message
}

func (mb *MessageBroker) Receive() (Message, bool) {
	select {
	case msg := <-mb.Channel:
		return msg, true
	default:
		return Message{}, false
	}
}

func (mb *MessageBroker) SendMsg(msg string, args ...any) {
	m := fmt.Sprintf(msg, args...)
	mb.Channel <- Message{Message: m}
}

func (mb *MessageBroker) SendError(msg string, err error) {
	mb.Channel <- Message{Message: msg, Err: err}
}

func (mb *MessageBroker) SendDone(msg, file string) {
	m := Message{Message: msg, Done: true, File: file}
	mb.Channel <- m
}

var ch = make(chan Message, 10)

func NewMessageBroker() *MessageBroker {
	mb := &MessageBroker{
		Channel: ch,
	}
	// fmt.Println(mb.Channel)
	return mb
}
