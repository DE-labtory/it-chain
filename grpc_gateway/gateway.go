package main

import "github.com/it-chain/it-chain-Engine/messaging"

type GrpcgateWay struct {
	messaing *messaging.Messaging
}

func (gw GrpcgateWay) Start() {
	gw.messaing.Consume()
}
