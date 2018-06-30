package icode

import "time"

type TransactionId = string

type Transaction struct {
	TxId          TransactionId
	PublishPeerId string
	TxHash        string
	TimeStamp     time.Time
	TxData        TxData
}

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
