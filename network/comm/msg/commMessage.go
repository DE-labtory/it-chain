package msg

import (
	"sync"
	"it-chain/network/protos"
)

type InnerMessage struct{
	Envelope *message.Envelope
	OnErr    func(error)
}

type OutterMessage struct{
	Envelope *message.Envelope
	Message *message.StreamMessage
	ConnectionID string
	Conn connection
	sync.Mutex
}

type connection interface {
	Send(envelope *message.Envelope, errCallBack func(error))
}

// Respond sends a msg to the source that sent the ReceivedMessageImpl
func (m *OutterMessage) Respond(envelope *message.Envelope, errCallBack func(error)) {

	m.Conn.Send(envelope, errCallBack)
}