package model

type TxpoolConfiguration struct {
	TimeoutMs          int64
	MaxTransactionByte int
}

func NewTxpoolConfiguration() TxpoolConfiguration {
	return TxpoolConfiguration{
		TimeoutMs:          1000,
		MaxTransactionByte: 1024,
	}
}
