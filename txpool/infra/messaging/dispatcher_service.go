package messaging

import (
	"github.com/it-chain/it-chain-Engine/txpool"
	"github.com/it-chain/midgard"
	"github.com/rs/xid"
	"errors"
)

type Publisher func(exchange string, topic string, data interface{}) (err error)

type MessageDispatcher struct {
	publisher Publisher // midgard.client
}

func NewDispatcher(publisher Publisher) *MessageDispatcher {
	return &MessageDispatcher{
		publisher: publisher,
	}
}


func (m MessageDispatcher) SendGrpcTransactions(transactions []*txpool.Transaction, leader txpool.Leader) error {
	if (len(transactions) == 0) {
		return errors.New("Empty transaction list proposed")
	}

	deliverCommand := txpool.MessageDeliverCommand{
		CommandModel: midgard.CommandModel{
			ID: xid.New().String(),
		},
		Transactions: transactions,
		Leader: leader,
	}

	return m.publisher("Command", "GrpcMessage", deliverCommand)
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


