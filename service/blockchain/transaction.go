package blockchain

import (
	"time"
	"it-chain/common"
)

type TransactionStatus int
type TxDataType string
type TransactionType int
type FunctionType string

const (

	status_TRANSACTION_UNCONFIRMED	Status	= 0
	status_TRANSACTION_CONFIRMED	Status	= 1
	status_TRANSACTION_UNKNOWN		Status	= 2

	invoke TxDataType = "invoke"
	query TxDataType = "query"

	general TransactionType = 0 + iota

	write = "write"
	read = "read"
	delete = "delete"
)

type Params struct {
	ParamsType	int
	Function 	FunctionType
	Args     	[]string
}

type TxData struct {
	Jsonrpc		string
	Method 		TxDataType
	Params 		Params
	ContractID	string
}

type Transaction struct {
	InvokePeerID      string
	TransactionID     string
	TransactionStatus Status
	TransactionType   TransactionType
	PublicKey         []byte
	Signature         []byte
	TransactionHash   string
	TimeStamp         time.Time
	TxData            *TxData
}

func CreateNewTransaction(peer_id string, tx_id string, status Status, tx_type TransactionType, key []byte, t time.Time, data *TxData) *Transaction{
	return &Transaction{InvokePeerID:peer_id, TransactionID:tx_id, TransactionStatus:status, TransactionType:tx_type, PublicKey:key, TimeStamp:t, TxData:data}
}

func MakeHashArg(tx Transaction) []string{
	sum := []string{tx.InvokePeerID, tx.TxData.Jsonrpc, string(tx.TxData.Method), string(tx.TxData.Params.Function), tx.TransactionID, tx.TimeStamp.String()}
	for _, str := range tx.TxData.Params.Args{ sum = append(sum, str) }
	return sum
}

func (tx *Transaction) GenerateHash() error{
	Arg := MakeHashArg(*tx)
	tx.TransactionHash = common.ComputeSHA256(Arg)
	return nil
}

func (tx Transaction) GenerateTransactionHash() string{
	Arg := MakeHashArg(tx)
	return common.ComputeSHA256(Arg)
}

func (tx *Transaction) GetTxHash() string{
	return tx.TransactionHash
}

func (tx Transaction) Validate() bool{
	if tx.GenerateTransactionHash() != tx.GetTxHash(){
		return false
	}
	return true
}