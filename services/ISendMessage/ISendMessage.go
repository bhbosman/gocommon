package ISendMessage

type ISendMessage interface {
	Send(message interface{}) error
}

type IMultiSendMessage interface {
	MultiSend(messages ...interface{})
}
