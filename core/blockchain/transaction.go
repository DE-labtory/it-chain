package blockchain

import (
	"time"
)

type TransactionStatus int

const (
	transactionUnconfirmed TransactionStatus = 0 + iota //unconfirmed block
	transactionConfirmed
)

type TransactionType int

const(
	general TransactionType = 0 + iota
)

type Transaction struct {
	invokePeerID string
	transactionID 	  string
	transactionStatus TransactionStatus
	transactionType   TransactionType
	publicKey         []byte
	signature         []byte
	transactionHash   string
	timeStamp         time.Time
}