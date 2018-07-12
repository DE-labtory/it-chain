package adapter

import (
	"errors"

	"github.com/it-chain/it-chain-Engine/txpool"
	"github.com/it-chain/it-chain-Engine/txpool/api"
)

var ErrNoEventID = errors.New("no event id ")

//////////////Event Handler

type EventHandler struct {
	transactionApi api.TransactionApi
}

func (e EventHandler) HandleBlockCommittedEvent(event txpool.BlockCommittedEvent) error {

	txs := event.Transactions

	for _, tx := range txs {
		err := e.transactionApi.DeleteTransaction(tx.TxId)

		if err != nil {
			return err
		}
	}

	return nil
}
