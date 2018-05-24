package messaging

import (
	"log"

	"time"

	"github.com/it-chain/it-chain-Engine/txpool"
	"github.com/it-chain/midgard"
	"github.com/rs/xid"
)

type TxCommandHandler struct {
	eventRepository midgard.Repository
	publisherId     string
}

func (t TxCommandHandler) HandleTxCreate(txCreateCommand txpool.TxCreateCommand) {

	events := make([]midgard.Event, 0)

	if txCreateCommand.GetID() == "" {
		log.Println("need id")
		return
	}

	id := string(txpool.TransactionId(xid.New().String()))
	timeStamp := time.Now()
	hash := txpool.CalTxHash(txCreateCommand.TxData, t.publisherId, txpool.TransactionId(id), txCreateCommand.TxType, timeStamp)

	events = append(events, txpool.TxCreatedEvent{
		EventModel: midgard.EventModel{
			ID: id,
		},
		PublishPeerId: t.publisherId,
		TxStatus:      txpool.VALID,
		TxType:        txpool.TxDataType(txCreateCommand.TxType),
		TxHash:        hash,
		TimeStamp:     timeStamp,
		TxData:        txCreateCommand.TxData,
	})

	err := t.eventRepository.Save(id, events...)

	if err != nil {
		log.Println(err.Error())
	}
}
