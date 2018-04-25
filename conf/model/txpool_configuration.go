package model

type TxpoolConfiguration struct {
	TimeoutMs          int
	MaxTransactionByte int
	RepositoryPath     string
}

func NewTxpoolConfiguration() TxpoolConfiguration {
	return TxpoolConfiguration{
		TimeoutMs:          100,
		MaxTransactionByte: 1024,
		RepositoryPath:     "empty",
	}
}
