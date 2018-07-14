package txpool

import (
	"time"

	"github.com/it-chain/midgard"
)

type TxCreatedEvent struct {
	midgard.EventModel
	PublishPeerId string
	TxStatus      int
	TxHash        string
	TimeStamp     time.Time
	Jsonrpc       string
	Method        string
	Params        Param
	ID            string
	ICodeID       string
}

func (tx TxCreatedEvent) GetTransaction() Transaction {

	return Transaction{
		TxId:          TransactionId(tx.EventModel.ID),
		PublishPeerId: tx.PublishPeerId,
		TxStatus:      TransactionStatus(tx.TxStatus),
		TxHash:        tx.TxHash,
		TxData: TxData{
			ICodeID: tx.ICodeID,
			Jsonrpc: tx.Jsonrpc,
			Method:  TxDataType(tx.Method),
			Params:  tx.Params,
			ID:      tx.ID,
		},
		TimeStamp: tx.TimeStamp,
	}
}

// when block committed check transaction and delete
type TxDeletedEvent struct {
	midgard.EventModel
}

type BlockCommittedEvent struct {
	midgard.EventModel
	Transactions []Transaction
}
