package blockchain

type MessageDispatcher interface {
	SendBlockCreatedEvent(block Block) error
}
