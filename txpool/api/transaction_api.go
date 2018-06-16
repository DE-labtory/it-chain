package api

import (
	"time"

	"github.com/it-chain/it-chain-Engine/txpool"
	"github.com/it-chain/midgard"
)

type TransactionApi struct {
	eventRepository *midgard.Repository
	publisherId     string
}

func NewTransactionApi(eventRepository *midgard.Repository, publisherId string) TransactionApi {
	return TransactionApi{
		publisherId:     publisherId,
		eventRepository: eventRepository,
	}
}

func (t TransactionApi) CreateTransaction(txID string, txData txpool.TxData) error {

	events := make([]midgard.Event, 0)

	timeStamp := time.Now()
	hash := txpool.CalTxHash(txData, t.publisherId, txpool.TransactionId(txID), timeStamp)

	events = append(events, txpool.TxCreatedEvent{
		EventModel: midgard.EventModel{
			ID:   txID,
			Type: "Transaction",
		},
		PublishPeerId: t.publisherId,
		TxStatus:      txpool.VALID,
		TxHash:        hash,
		TimeStamp:     timeStamp,
		TxData:        txData,
	})

	err := t.eventRepository.Save(txID, events...)

	if err != nil {
		return err
	}

	return nil
}
