package transaction

type TxData struct {
	Jsonrpc string
	Method  string
	Params  []string
	ICodeId string `json:"id"`
}

func NewTxData(jsonrpc string, method string, params []string, iCodeId string) *TxData{
	return &TxData{
		Jsonrpc: jsonrpc,
		Method:  method,
		Params:  params,
		ICodeId: iCodeId,
	}
}