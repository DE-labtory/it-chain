package txpool

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

func NewTxData(jsonrpc string, method TxDataType, params Param, iCodeId string, id string) *TxData {
	return &TxData{
		Jsonrpc: jsonrpc,
		Method:  method,
		Params:  params,
		ID:      id,
		ICodeID: iCodeId,
	}
}
