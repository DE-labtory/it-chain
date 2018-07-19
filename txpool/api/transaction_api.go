package api

import (
	"log"

	"github.com/it-chain/engine/core/eventstore"
	"github.com/it-chain/engine/txpool"
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

	log.Printf("create transaction: [%v]", txData)

	tx, err := txpool.CreateTransaction(t.publisherId, txData)

	if err != nil {
		log.Printf("fail to transaction: [%v]", err)
		return tx, err
	}

	log.Printf("transaction created: [%v]", tx)
	return tx, nil
}

func (t TransactionApi) DeleteTransaction(id txpool.TransactionId) error {

	log.Printf("delete transaction: [%v]", id)

	tx := &txpool.Transaction{}

	if err := eventstore.Load(tx, id); err != nil {
		log.Printf("fail to delete transaction: [%v]", id)
		return err
	}

	return txpool.DeleteTransaction(*tx)
}
