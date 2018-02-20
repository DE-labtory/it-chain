package domain

import (
	"time"
	"it-chain/common"
	"errors"
	pb "it-chain/network/protos"
)


type TransactionStatus int
type TxDataType string
type TransactionType int
type FunctionType string

const (

	Status_TRANSACTION_UNCONFIRMED	TransactionStatus	= 0
	Status_TRANSACTION_CONFIRMED	TransactionStatus	= 1
	Status_TRANSACTION_UNKNOWN		TransactionStatus	= 2

	Invoke TxDataType = "invoke"
	Query TxDataType = "query"

	General TransactionType = 0 + iota

	Write = "write"
	Read = "read"
	Delete = "delete"
)

type Params struct {
	ParamsType	int
	Function 	string
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
	TransactionStatus TransactionStatus
	TransactionType   TransactionType
	PublicKey         []byte
	Signature         []byte
	TransactionHash   string
	TimeStamp         time.Time
	TxData            *TxData
}

func CreateNewTransaction(peer_id string, tx_id string, tx_type TransactionType, t time.Time, data *TxData) *Transaction{

	transaction := &Transaction{InvokePeerID:peer_id, TransactionID:tx_id, TransactionStatus:Status_TRANSACTION_UNKNOWN, TransactionType:tx_type, TimeStamp:t, TxData:data}
	transaction.GenerateHash()
	return transaction
}

func SetTxMethodParameters(params_type int, function string, args []string) Params{
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

func (tx *Transaction) SignHash() (ret bool, err error){
	signature := []byte("temp")
	tx.PublicKey = []byte("temp")

	if signature != nil {
		tx.Signature = signature
		err = nil
		ret = true
	} else {
		err = errors.New("transaction signature fail")
		ret = false
	}
	return ret, err
}

//todo test
func FromProtoTxData(ptxData pb.TxData) *TxData{

	txData := &TxData{
		Params: Params{
			Args:ptxData.Params.Args,
			Function: ptxData.Params.Function,
			ParamsType: int(ptxData.Params.ParamsType),
		},
		ContractID:ptxData.ContractID,
		Jsonrpc: ptxData.Jsonrpc,
	}

	if ptxData.Method == pb.TxData_Invoke{
		txData.Method = Invoke
	}

	if ptxData.Method == pb.TxData_Query{
		txData.Method = Query
	}

	return txData
}

//todo test
func ToProtoTxData(t TxData) *pb.TxData{

	txData := &pb.TxData{
		Params: &pb.Params{
			Args:t.Params.Args,
			Function: t.Params.Function,
			ParamsType: int32(t.Params.ParamsType),
		},
		ContractID:t.ContractID,
		Jsonrpc: t.Jsonrpc,
	}

	if t.Method == Invoke{
		txData.Method =  pb.TxData_Invoke
	}

	if t.Method == Query{
		txData.Method = pb.TxData_Invoke
	}

	return txData
}

//todo test
func ToProtoTransaction(t Transaction) *pb.Transaction{
	transaction := &pb.Transaction{
		TransactionHash: t.TransactionHash,
		TransactionStatus: pb.Transaction_Status(t.TransactionStatus),
		TransactionID: t.TransactionID,
		InvokePeerID: t.InvokePeerID,
		TxData: ToProtoTxData(*t.TxData),
	}

	return transaction
}
