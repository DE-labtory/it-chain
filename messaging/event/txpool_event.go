package event

type TransactionReceiveEvent struct {
	PeerId      string
	Transaction []byte
}

type TransactionSendEvent struct {
	LeaderId    string
	Transaction []byte
}

type BlockProposeEvent struct {
	TransactionList []byte
}

type TransactionCreateEvent struct {
	TransactionData []byte
}
