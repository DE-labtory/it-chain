package transaction

type TxData struct {
	Jsonrpc string
	Method  string
	Params  []string
	ICodeId string `json:"id"`
}
