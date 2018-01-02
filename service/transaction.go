package service

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
	paramsType int
	function   FunctionType
	args       []string
}

type TxData struct {
	jsonrpc string
	method  TxDataType
	params  Params
	contractID string
}

type Transaction struct {
	invokePeerID      string
	transactionID     string
	transactionStatus TransactionStatus
	transactionType   TransactionType
	publicKey         []byte
	signature         []byte
	transactionHash   string
	timeStamp         time.Time
	txData            TxData
}