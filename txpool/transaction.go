package txpool

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/it-chain/it-chain-Engine/common"
	"github.com/it-chain/it-chain-Engine/core/eventstore"
	"github.com/it-chain/midgard"
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

//Aggregate root must implement aggregate interface
type Transaction struct {
	TxId          TransactionId
	PublishPeerId string
	TxStatus      TransactionStatus
	TxHash        string
	TimeStamp     time.Time
	TxData        TxData
}

// must implement id method
func (t Transaction) GetID() string {
	return string(t.TxId)
}

// must implement on method
func (t *Transaction) On(event midgard.Event) error {

	switch v := event.(type) {

	case *TxCreatedEvent:
		t.TxId = TransactionId(v.ID)
		t.PublishPeerId = v.PublishPeerId
		t.TxStatus = v.TxStatus
		t.TxHash = v.TxHash
		t.TimeStamp = v.TimeStamp
		t.TxData = v.TxData

	default:
		return errors.New(fmt.Sprintf("unhandled event [%s]", v))
	}

	return nil
}

func (t Transaction) Serialize() ([]byte, error) {
	return common.Serialize(t)
}

func Deserialize(b []byte, transaction *Transaction) error {

	err := json.Unmarshal(b, transaction)

	if err != nil {
		return err
	}

	return nil
}

func CalTxHash(txData TxData, publishPeerId string, txId TransactionId, timeStamp time.Time) string {

	hashArgs := []string{
		txData.Jsonrpc,
		string(txData.Method),
		string(txData.Params.Function),
		txData.ICodeID,
		publishPeerId,
		timeStamp.String(),
		string(txId),
	}

	for _, str := range txData.Params.Args {
		hashArgs = append(hashArgs, str)
	}

	return common.ComputeSHA256(hashArgs)
}

func CreateTransaction(publisherId string, txData TxData) (Transaction, error) {

	id := xid.New().String()
	timeStamp := time.Now()
	hash := CalTxHash(txData, publisherId, TransactionId(id), timeStamp)

	event := &TxCreatedEvent{
		EventModel: midgard.EventModel{
			ID:   id,
			Type: "Transaction",
		},
		PublishPeerId: publisherId,
		TxStatus:      VALID,
		TxHash:        hash,
		TimeStamp:     timeStamp,
		TxData:        txData,
	}

	tx := Transaction{}
	tx.On(event)

	if err := eventstore.Save(tx.GetID(), event); err != nil {
		return tx, err
	}

	return tx, nil
}

//TxData Declaration
const (
	Invoke TxDataType = "invoke"
	Query  TxDataType = "query"
)

type TxDataType string

type TxData struct {
	Jsonrpc string
	Method  TxDataType
	Params  Param
	ID      string
	ICodeID string
}

type Param struct {
	Function string
	Args     []string
}

func NewTxData(jsonrpc string, method TxDataType, params Param, iCodeId string, id string) *TxData {
	return &TxData{
		Jsonrpc: jsonrpc,
		Method:  method,
		Params:  params,
		ID:      id,
		ICodeID: iCodeId,
	}
}

//Transaction Repository interface
type TransactionRepository interface {
	Save(transaction Transaction) error
	Remove(id TransactionId) error
	FindById(id TransactionId) (*Transaction, error)
	FindAll() ([]*Transaction, error)
}
