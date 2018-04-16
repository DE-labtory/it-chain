package factory

import (
	"time"
	"github.com/it-chain/it-chain-Engine/txpool/domain/model/transaction"
	"github.com/rs/xid"
)

//todo peer id 타입 생기면 'publishPeerId' string->해당 타입으로 변경
func CreateNewTransaction(publishPeerId string, txType transaction.TxType, txData *transaction.TxData) *transaction.Transaction{
	return &transaction.Transaction{
		TxId:          transaction.TransactionId(xid.New().String()),
		PublishPeerId: publishPeerId,
		TxStatus:      transaction.STATUS_TX_UNKNOWN,
		TxType:        txType,
		TxHash:        transaction.HASH_NEED_CALC,
		TimeStamp:     time.Now(),
		TxData:        txData,
	}
}