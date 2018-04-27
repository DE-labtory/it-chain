package main

import "github.com/it-chain/it-chain-Engine/messaging"

type GrpcgateWay struct {
	messaing *messaging.Messaging
	ec       EventConsumer
}

func (gw GrpcgateWay) Start() {
}
