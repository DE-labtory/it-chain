package api

import (
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

	pool := txpool.LoadTransactionPool()

	if err = pool.Add(tx); err != nil {

		return tx, err
	}

	return tx, nil
}
