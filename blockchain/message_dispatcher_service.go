package blockchain

type MessageDispatcher interface {
	SendBlockValidateCommand(block Block) error
}
