package transaction

import (
	"time"
	"github.com/rs/xid"
)

const (
	//Tx 상태 Const
	STATUS_TX_UNCONFIRM TxStatus = iota
	STATUS_TX_CONFIRM
	STATUS_TX_INVALID

	//Tx 타입 Const
	TYPE_DEPLOY TxType = iota
	TYPE_INVOKE
	TYPE_QUERY

	HASH_NEED_CALC = "HASH_NEED_CALC"
)

type TransactionId string
type TxStatus int
type TxType int

type Transaction struct {
	TxId          TransactionId
	PublishPeerId string
	TxStatus      TxStatus
	TxType        TxType
	TxHash        string
	TimeStamp     time.Time
	TxData        *TxData
}

func NewTransaction(publishPeerId string, txType TxType, txData *TxData) *Transaction{
	return &Transaction{
		TxId:          TransactionId(xid.New().String()),
		PublishPeerId: publishPeerId,
		TxStatus:      STATUS_TX_UNCONFIRM,
		TxType:        txType,
		TxHash:        HASH_NEED_CALC,
		TimeStamp:     time.Now(),
		TxData:        txData,
	}
}

