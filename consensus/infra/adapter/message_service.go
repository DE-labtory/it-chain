package adapter

import (
	c "github.com/it-chain/it-chain-Engine/consensus"
)

type Publisher func(exchange string, topic string, data interface{}) (err error)

type MessageService struct {
	publisher Publisher
}

func NewMessageService(publisher Publisher) *MessageService {
	return &MessageService{
		publisher: publisher,
	}
}

// todo
func (m MessageService) BroadcastMsg(Msg c.Serializable, representatives []*c.Representative) {

}

// todo
func (m MessageService) CreateConfirmedBlock(block c.ProposedBlock) {

}
