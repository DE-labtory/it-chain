package blockchain

import (
	"time"
	"it-chain/common"
	"strings"
)

type TransactionStatus int
type TxDataType string
type TransactionType int
type FunctionType string

const (
	invoke TxDataType = "invoke"
	query TxDataType = "query"

	general TransactionType = 0 + iota

	write = "write"
	read = "read"
	delete = "delete"
)

type Params struct {
	paramsType	int
	function 	FunctionType
	args     	[]string
}

type TxData struct {
	jsonrpc		string
	method 		TxDataType
	params 		Params
	contractID	string
}

type Transaction struct {
	invokePeerID      string
	transactionID     string
	transactionStatus Status
	transactionType   TransactionType
	publicKey         []byte
	signature         []byte
	transactionHash   [32]uint8
	timeStamp         time.Time
	txData            *TxData
}

func CreateNewTransaction(peer_id string, tx_id string, status Status, tx_type TransactionType, key []byte, hash [32]uint8, t time.Time, data *TxData) *Transaction{
	return &Transaction{invokePeerID:peer_id, transactionID:tx_id, transactionStatus:status, transactionType:tx_type, publicKey:key, transactionHash:hash, timeStamp:t, txData:data}
}

func MakeHashArg(tx Transaction) []byte{
	sum := []string{tx.invokePeerID, tx.txData.jsonrpc, string(tx.txData.method), string(tx.txData.params.function), tx.transactionID, tx.timeStamp.String()}
	for _, str := range tx.txData.params.args{ sum = append(sum, str) }
	str := strings.Join(sum, ",")
	return []byte(str)
}

func (tx *Transaction) GenerateHash() error{
	Arg := MakeHashArg(*tx)
	tx.transactionHash = common.ComputeSHA256(Arg)
	return nil
}

func (tx Transaction) GenerateTransactionHash() [32]uint8{
	Arg := MakeHashArg(tx)
	return common.ComputeSHA256(Arg)
}

func (tx *Transaction) GetTxHash() [32]uint8{
	return tx.transactionHash
}

func (tx Transaction) Validate() bool{
	if tx.GenerateTransactionHash() != tx.GetTxHash(){
		return false
	}
	return true
}