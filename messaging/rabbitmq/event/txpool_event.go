package event

type TransactionReceiveEvent struct {
	PeerId      string
	Transaction []byte
}

//todo issue #165 에 따른 변경필요
type TransactionSendEvent struct {
	LeaderId    string
	Transaction []byte
}

type BlockProposeEvent struct {
	TransactionList []byte
}

type TransactionCreateEvent struct {
	PeerId			string
	TxDataType		string
	TransactionData []byte
}
