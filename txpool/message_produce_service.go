package txpool

import (
	"github.com/it-chain/it-chain-Engine/txpool/domain/model"
	"github.com/it-chain/it-chain-Engine/txpool/domain/model/transaction"
)

type Publish func(topic string, data []byte) error

type MessageProducer interface {
	SendTransactions(transactions []transaction.Transaction, leader model.Leader) error
	ProposeBlock(transactions []transaction.Transaction) error
}
