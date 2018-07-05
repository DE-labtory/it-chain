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
	TxData        TxData
}

func (tx TxCreatedEvent) GetTransaction() Transaction {

	return Transaction{
		TxId:          TransactionId(tx.ID),
		PublishPeerId: tx.PublishPeerId,
		TxStatus:      tx.TxStatus,
		TxHash:        tx.TxHash,
		TxData:        tx.TxData,
	}
}

type TxAddedPoolEvent struct {
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
}

type TxDeletedFromPoolEvent struct {
	midgard.EventModel
}

type TxCommitedToStageEvent struct {
	midgard.EventModel
}

type LeaderChangedEvent struct {
	midgard.EventModel
}
