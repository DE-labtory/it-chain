package api

import (
	"log"
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

func (t TransactionApi) CreateTransaction(txCreateCommand txpool.TxCreateCommand) {

	events := make([]midgard.Event, 0)

	if txCreateCommand.GetID() == "" {
		log.Println("need id")
		return
	}

	id := txCreateCommand.GetID()
	timeStamp := time.Now()
	hash := txpool.CalTxHash(txCreateCommand.TxData, t.publisherId, txpool.TransactionId(id), timeStamp)

	events = append(events, txpool.TxCreatedEvent{
		EventModel: midgard.EventModel{
			ID:   id,
			Type: "Transaction",
		},
		PublishPeerId: t.publisherId,
		TxStatus:      txpool.VALID,
		TxHash:        hash,
		TimeStamp:     timeStamp,
		TxData:        txCreateCommand.TxData,
	})

	err := t.eventRepository.Save(id, events...)

	if err != nil {
		log.Println(err.Error())
	}
}
