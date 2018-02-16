package domain

import (
	"time"
	"it-chain/common"
)

type TransactionStatus int
type TxDataType string
type TransactionType int
type FunctionType string

const (

	Status_TRANSACTION_UNCONFIRMED	Status	= 0
	Status_TRANSACTION_CONFIRMED	Status	= 1
	Status_TRANSACTION_UNKNOWN		Status	= 2

	Invoke TxDataType = "invoke"
	Query TxDataType = "query"

	General TransactionType = 0 + iota

	Write = "write"
	Read = "read"
	Delete = "delete"
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

func CreateNewTransaction(peer_id string, tx_id string, tx_type TransactionType, t time.Time, data *TxData) *Transaction{
	return &Transaction{InvokePeerID:peer_id, TransactionID:tx_id, TransactionStatus:Status_TRANSACTION_UNKNOWN, TransactionType:tx_type, TimeStamp:t, TxData:data}
}

func SetTxMethodParameters(params_type int, function FunctionType, args []string) Params{
	return Params{params_type, function, args}
}

func SetTxData(jsonrpc string, method TxDataType, params Params, contract_id string) *TxData{
	return &TxData{jsonrpc, method, params, contract_id}
}

func MakeHashArg(tx Transaction) []string{
	sum := []string{tx.InvokePeerID, tx.TxData.Jsonrpc, string(tx.TxData.Method), string(tx.TxData.Params.Function), tx.TransactionID, tx.TimeStamp.String()}
	for _, str := range tx.TxData.Params.Args{ sum = append(sum, str) }
	return sum
}

func (tx *Transaction) GenerateHash() {
	Arg := MakeHashArg(*tx)
	tx.TransactionHash = common.ComputeSHA256(Arg)
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