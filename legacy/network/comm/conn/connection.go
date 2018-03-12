package conn

import "github.com/it-chain/it-chain-Engine/legacy/network/protos"

type Connection interface{
	Send(envelope *message.Envelope, successCallBack func(interface{}),errCallBack func(error))
	Close()
}