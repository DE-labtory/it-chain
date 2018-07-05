package txpool

import (
	"time"

	"github.com/it-chain/midgard"
)

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
