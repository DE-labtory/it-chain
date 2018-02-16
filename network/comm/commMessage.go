package comm

import (
	"sync"
	"it-chain/network/protos"
)

type innerMessage struct{
	envelope *message.Envelope
	onErr    func(error)
}

type OutterMessage struct{
	Envelope *message.Envelope
	Message *message.StreamMessage
	ConnectionID string
	Conn Connection
	sync.Mutex
}

// Respond sends a msg to the source that sent the ReceivedMessageImpl
func (m *OutterMessage) Respond(envelope *message.Envelope, errCallBack func(error)) {

	m.Conn.Send(envelope, errCallBack)
}