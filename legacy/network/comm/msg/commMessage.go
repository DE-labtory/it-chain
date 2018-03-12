package msg

import (
	"sync"
	"github.com/it-chain/it-chain-Engine/legacy/network/protos"
)

type InnerMessage struct{
	Envelope *message.Envelope
	OnErr    func(error)
	OnSuccess func(interface{})
}

type OutterMessage struct{
	Envelope *message.Envelope
	Message *message.StreamMessage
	ConnectionID string
	Conn connection
	sync.Mutex
}

type connection interface {
	Send(envelope *message.Envelope, successCallBack func(interface{}), errCallBack func(error))
}

// Respond sends a msg to the source that sent the ReceivedMessageImpl
func (m *OutterMessage) Respond(envelope *message.Envelope, successCallBack func(interface{}), errCallBack func(error)) {

	m.Conn.Send(envelope, successCallBack, errCallBack)
}