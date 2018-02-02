package comm

import "it-chain/network/protos"

type Connection interface{
	Send(envelope *message.Envelope, errCallBack func(error))
	Close()
}