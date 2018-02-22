package conn

import "it-chain/network/protos"

type Connection interface{
	Send(envelope *message.Envelope, successCallBack func(interface{}),errCallBack func(error))
	Close()
}