package rabbitmq

import (
	"github.com/it-chain/it-chain-Engine/txpool/domain/model"
	"github.com/it-chain/it-chain-Engine/txpool/domain/model/transaction"
	"github.com/it-chain/it-chain-Engine/txpool/domain/service"
)

// todo need test code
type MessageProducer struct {
	Publish service.Publish
}

func NewMessageProducer(publish service.Publish) *MessageProducer {
	return &MessageProducer{
		Publish: publish,
	}
}

// todo impl SendTransactions
func (mp *MessageProducer) SendTransactions(transaction []transaction.Transaction, leader model.Leader) error {
	panic("implement me")
}

// todo impl ProposeBlock
func (mp *MessageProducer) ProposeBlock(transactions []transaction.Transaction) error {
	panic("implement me")
}
