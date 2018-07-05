package api

import (
	"github.com/it-chain/it-chain-Engine/txpool"
	"github.com/it-chain/midgard"
)

type TransactionApi struct {
	eventRepository midgard.EventRepository
	publisherId     string
}

func NewTransactionApi(eventRepository midgard.EventRepository, publisherId string) TransactionApi {
	return TransactionApi{
		publisherId:     publisherId,
		eventRepository: eventRepository,
	}
}

func (t TransactionApi) CreateTransaction(txData txpool.TxData) (txpool.Transaction, error) {

	tx, err := txpool.CreateTransaction(t.publisherId, txData)

	if err != nil {
		return tx, err
	}

	return tx, nil
}
