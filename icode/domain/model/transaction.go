package model

import "time"

type Transaction struct {
	InvokePeerID      string
	TransactionID     string
	TransactionStatus TransactionStatus
	TransactionType   TransactionType
	TransactionHash   string
	TimeStamp         time.Time
	TxData            *TxData
}
