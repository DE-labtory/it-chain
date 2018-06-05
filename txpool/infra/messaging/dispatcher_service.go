package messaging

import (
	"github.com/it-chain/it-chain-Engine/txpool"
	"github.com/it-chain/midgard"
	"github.com/rs/xid"
	"errors"
)

type Publisher func(exchange string, topic string, data interface{}) (err error) 	//해당 publish함수는 midgard 에서 의존성 주입을 받기 위해 interface로 작성한다.
																						//모든 의존성 주입은 컴포넌트.go 에서 이루어짐

type MessageDispatcher struct {
	publisher Publisher // midgard.client
}

func NewDispatcher(publisher Publisher) *MessageDispatcher {
	return &MessageDispatcher{
		publisher: publisher,
	}
}


func (m MessageDispatcher) SendLeaderTransactions(transactions []*txpool.Transaction, leader txpool.Leader) error {
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


