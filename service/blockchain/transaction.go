package blockchain

import (
	"time"
)

type TransactionStatus int
type TxDataType string
type TransactionType int
type FunctionType string

const (

	transactionUnconfirmed TransactionStatus = 0 + iota //unconfirmed block
	transactionConfirmed

	invoke TxDataType = "invoke"
	query TxDataType = "query"

	general TransactionType = 0 + iota

	write = "write"
	read = "read"
	delete = "delete"
)

type Params struct {
	ParamsType int
	Function   FunctionType
	Args       []string
}

type TxData struct {
	Jsonrpc string
	Method  TxDataType
	Params  Params
	ContractID string
}

type Transaction struct {
	InvokePeerID      string
	TransactionID     string
	TransactionStatus TransactionStatus
	TransactionType   TransactionType
	PublicKey         []byte
	Signature         []byte
	TransactionHash   string
	TimeStamp         time.Time
	TxData            TxData
}