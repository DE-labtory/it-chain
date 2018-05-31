package messaging

import (
	"github.com/it-chain/it-chain-Engine/txpool"
	"github.com/it-chain/midgard"
	"github.com/rs/xid"
	"errors"
)

type Publisher func(exchange string, topic string, data interface{}) (err error)



//todo implement create command using transaction and leader and send to rabbitmq
type MessageDispatcher struct {
	publisher Publisher // midgard.client
}

func NewDispatcher(publisher Publisher) *MessageDispatcher {
	return &MessageDispatcher{
		publisher: publisher,
	}
}


//todo implement sendTransactionCommand 정의 해야함
func (m MessageDispatcher) SendTransactions(transactions []*txpool.Transaction, leader txpool.Leader) error {
	if (len(transactions) == 0) {
		return errors.New("Empty transaction list proposed")
	}

	deliverCommand := txpool.SendTransactionsCommand{
		CommandModel: midgard.CommandModel{
			ID: xid.New().String(),
		},
		Transactions: transactions,
		Leader: leader,
	}

	return m.publisher("Command", "Transactions", deliverCommand)
}

func (m MessageDispatcher) ProposeBlock(transactions []txpool.Transaction) error {
	if (len(transactions) == 0) {
		return errors.New("Empty transaction list proposed")
	}

	deliverCommand := txpool.ProposeBlockCommand{
		CommandModel: midgard.CommandModel{
			ID: xid.New().String(),
		},
		Transactions: transactions,
	}

	return m.publisher("Command", "Block", deliverCommand)
}


