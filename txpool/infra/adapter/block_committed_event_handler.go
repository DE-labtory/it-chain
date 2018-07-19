package adapter

import (
	"errors"

	"github.com/it-chain/engine/txpool"
	"github.com/it-chain/engine/txpool/api"
)

var ErrNoEventID = errors.New("no event id ")

type BlockCommittedEventHandler struct {
	transactionApi api.TransactionApi
}

func (e BlockCommittedEventHandler) HandleBlockCommittedEvent(event txpool.BlockCommittedEvent) error {

	txs := event.Transactions

	for _, tx := range txs {
		err := e.transactionApi.DeleteTransaction(tx.TxId)

		if err != nil {
			return err
		}
	}

	return nil
}
