package txpool

import (
	"time"

	"encoding/json"

	"github.com/it-chain/it-chain-Engine/common"
	"github.com/rs/xid"
)

const (
	VALID   TransactionStatus = 0
	INVALID TransactionStatus = 1

	General TransactionType = 0 + iota
)

type TransactionId string

func (tId TransactionId) ToString() string {
	return string(tId)
}

type TransactionStatus int
type TransactionType int
type Transaction struct {
	TxId          TransactionId
	PublishPeerId string
	TxStatus      TransactionStatus
	TxType        TxDataType
	TxHash        string
	TimeStamp     time.Time
	TxData        *TxData
}

func NewTransaction(publishPeerId string, txType TxDataType, txData *TxData) *Transaction {
	tx := Transaction{
		TxId:          TransactionId(xid.New().String()),
		PublishPeerId: publishPeerId,
		TxStatus:      VALID,
		TxType:        txType,
		TxHash:        "",
		TimeStamp:     time.Now(),
		TxData:        txData,
	}

	tx.TxHash = tx.CalcHash()
	return &tx
}

func (t Transaction) Serialize() ([]byte, error) {
	return common.Serialize(t)
}

func (t Transaction) GetID() string {
	return string(t.TxId)
}

func Deserialize(b []byte, transaction *Transaction) error {
	err := json.Unmarshal(b, transaction)

	if err != nil {
		return err
	}

	return nil
}

func (t Transaction) CalcHash() string {
	hashArgs := []string{t.TxData.Jsonrpc, string(t.TxData.Method), string(t.TxData.Params.Function), t.TxData.ICodeID, t.PublishPeerId, t.TimeStamp.String(), string(t.TxId), string(t.TxType)}
	for _, str := range t.TxData.Params.Args {
		hashArgs = append(hashArgs, str)
	}
	return common.ComputeSHA256(hashArgs)
}

type TransactionRepository interface {
	Save(transaction Transaction) error
	Remove(id TransactionId) error
	FindById(id TransactionId) (*Transaction, error)
	FindAll() ([]*Transaction, error)
}
