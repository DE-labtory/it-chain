package event

// todo define event
type TransactionReceiveEvent struct {
	PeerId      string
	Transaction []byte
}

// todo define event
type TransactionSendEvent struct {
	LeaderId    string
	Transaction []byte
}

// todo define event
type BlockProposeEvent struct {
	TransactionList []byte
}

type TransactionCreateEvent struct {
	TransactionData []byte
}
