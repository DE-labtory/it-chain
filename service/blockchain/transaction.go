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
	txData            TxData
}

func (tx *Transaction) init(peer_id string, tx_id string, status Status, tx_type TransactionType, key []byte, sign []byte, hash [32]uint8, t time.Time, data *TxData){
	tx.invokePeerID 		= peer_id
	tx.transactionID 		= tx_id
	tx.transactionStatus	= status
	tx.transactionType		= tx_type
	tx.publicKey			= key
	tx.signature			= sign
	tx.transactionHash		= hash
	tx.timeStamp			= t
	tx.txData				= *data
}

func (tx *Transaction) GenerateHash() error{
	sum_string := []string{tx.invokePeerID, tx.txData.jsonrpc, string(tx.txData.method), string(tx.txData.params.function)}
	for _, str := range tx.txData.params.args{
		sum_string = append(sum_string, str)
	}
	str := strings.Join(sum_string, ",")
	tx.transactionHash = common.ComputeSHA256([]byte(str))
	return nil
}

func (tx Transaction) GenerateTransactionHash() [32]uint8{
	sum_string := []string{tx.invokePeerID, tx.txData.jsonrpc, string(tx.txData.method), string(tx.txData.params.function)}
	for _, str := range tx.txData.params.args{
		sum_string = append(sum_string, str)
	}
	str := strings.Join(sum_string, ",")
	return common.ComputeSHA256([]byte(str))
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