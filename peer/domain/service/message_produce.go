package service

import "github.com/it-chain/it-chain-Engine/peer/domain/model"

type Publish func(topic string, data []byte) error

type MessageProducer interface {
	RequestLeaderInfo(peer model.Peer) error
	DeliverLeaderInfo(toPeer model.Peer, leader model.Peer) error
	LeaderUpdateEvent(leader model.Peer) error
}
