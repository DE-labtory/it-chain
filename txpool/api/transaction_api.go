package api

import (
	"log"

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

	log.Printf("create transaction: [%s]", txData)

	tx, err := txpool.CreateTransaction(t.publisherId, txData)

	if err != nil {
		log.Printf("fail to transaction: [%s]", err)
		return tx, err
	}

	log.Printf("transaction created: [%s]", tx)
	return tx, nil
}

func (t TransactionApi) DeleteTransaction(id txpool.TransactionId) error {

	log.Printf("delete transaction: [%s]", id)

	tx := &txpool.Transaction{}

	if err := eventstore.Load(tx, id); err != nil {
		log.Printf("fail to delete transaction: [%s]", id)
		return err
	}

	return txpool.DeleteTransaction(*tx)
}
