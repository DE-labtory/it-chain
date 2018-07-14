package icode

import "time"

const (
	Ready = iota
	ExcuteSuccess
	ExcuteFail
)

const (
	Invoke TxDataType = "invoke"
	Query  TxDataType = "query"
)

type TransactionId = string
type TxStatus = int
type TxDataType = string

type Transaction struct {
	TxId      TransactionId
	TimeStamp time.Time
	TxData    TxData
}

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
