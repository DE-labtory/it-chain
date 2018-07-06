package txpool

import (
	"time"

	"github.com/it-chain/midgard"
)

type TxCreatedEvent struct {
	midgard.EventModel
	PublishPeerId string
	TxStatus      TransactionStatus
	TxHash        string
	TimeStamp     time.Time
	Jsonrpc       string
	Method        TxDataType
	Params        Param
	ID            string
	ICodeID       string
	Stage         TransactionStage
}

func (tx TxCreatedEvent) GetTransaction() Transaction {

	return Transaction{
		TxId:          TransactionId(tx.ID),
		PublishPeerId: tx.PublishPeerId,
		TxStatus:      tx.TxStatus,
		TxHash:        tx.TxHash,
		TxData: TxData{
			ICodeID: tx.ICodeID,
			Jsonrpc: tx.Jsonrpc,
			Method:  tx.Method,
			Params:  tx.Params,
			ID:      tx.ID,
		},
		TimeStamp: tx.TimeStamp,
	}
}

type TxStagedEvent struct {
	midgard.EventModel
}

// when block
type TxDeleteEvent struct {
	midgard.EventModel
}

type LeaderChangedEvent struct {
	midgard.EventModel
}
