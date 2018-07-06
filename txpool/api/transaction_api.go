package api

import (
	"github.com/it-chain/it-chain-Engine/core/eventstore"
	"github.com/it-chain/it-chain-Engine/txpool"
)

type TransactionApi struct {
	publisherId string
}

func NewTransactionApi(publisherId string) TransactionApi {
	return TransactionApi{
		publisherId: publisherId,
	}
}

func (t TransactionApi) CreateTransaction(txData txpool.TxData) (txpool.Transaction, error) {

	tx, err := txpool.CreateTransaction(t.publisherId, txData)

	if err != nil {
		return tx, err
	}

	return tx, nil
}

func (t TransactionApi) DeleteTransaction(id txpool.TransactionId) error {

	tx := &txpool.Transaction{}

	if err := eventstore.Load(tx, id); err != nil {
		return err
	}

	return txpool.DeleteTransaction(*tx)
}
