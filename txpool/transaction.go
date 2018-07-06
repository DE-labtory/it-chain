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

type TransactionId = string

type TransactionStatus int
type TransactionType int

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
		t.TxData = TxData{
			ID:      v.ID,
			Params:  v.Params,
			Method:  v.Method,
			Jsonrpc: v.Jsonrpc,
			ICodeID: v.ICodeID,
		}

	case *TxDeletedEvent:
		t.TxId = ""

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
			Type: "transaction.created",
		},
		PublishPeerId: publisherId,
		TxStatus:      VALID,
		TxHash:        hash,
		TimeStamp:     timeStamp,
		ID:            txData.ID,
		ICodeID:       txData.ICodeID,
		Jsonrpc:       txData.Jsonrpc,
		Method:        txData.Method,
		Params:        txData.Params,
	}

	tx := &Transaction{}

	if err := saveAndOn(tx, event); err != nil {
		return *tx, err
	}

	return *tx, nil
}

func DeleteTransaction(transaction Transaction) error {

	event := &TxDeletedEvent{
		EventModel: midgard.EventModel{
			ID: transaction.TxId,
		},
	}

	tx := &Transaction{}

	if err := saveAndOn(tx, event); err != nil {
		return err
	}

	return nil
}

//apply on aggrgate and publish to eventstore
func saveAndOn(aggregate midgard.Aggregate, event midgard.Event) error {

	//must do call on func first!!!
	//after save events if aggregate.On failed then data inconsistency will be occurred
	if err := aggregate.On(event); err != nil {
		return err
	}

	if err := eventstore.Save(event.GetID(), event); err != nil {
		return err
	}

	return nil
}
