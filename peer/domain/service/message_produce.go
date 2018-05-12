package service

import "github.com/it-chain/it-chain-Engine/peer/domain/model"

type Publish func(topic string, data []byte) error

type MessageProducer interface {
	RequestLeaderInfo(ipAddress string) error
	LeaderUpdateEvent(leader model.Peer) error
}
