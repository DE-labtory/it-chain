package blockchain

type Result struct {
	TxId    TransactionId
	Data    map[key]value
	Success bool
}

type key = string
type value = string
