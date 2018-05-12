package rabiitmq

import "github.com/it-chain/it-chain-Engine/peer/domain/service"

type MessageProducer struct {
	Publish service.Publish
}

func NewMessageProducer(publish service.Publish) *MessageProducer {
	return &MessageProducer{
		Publish: publish,
	}
}
