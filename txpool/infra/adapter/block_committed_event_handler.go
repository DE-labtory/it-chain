package adapter

import (
	"errors"

	"github.com/it-chain/it-chain-Engine/txpool"
	"github.com/it-chain/it-chain-Engine/txpool/api"
)

var ErrNoEventID = errors.New("no event id ")

//////////////Event Handler

type EventHandler struct {
	transacionApi api.TransactionApi
}

func (e EventHandler) HandleBlockCommitedEvent(event txpool.BlockCommittedEvent) error {

	txs := event.Transactions

	for _, tx := range txs {
		err := e.transacionApi.DeleteTransaction(tx.TxId)

		if err != nil {
			return err
		}
	}

	return nil
}
