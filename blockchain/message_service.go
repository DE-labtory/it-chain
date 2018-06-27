package blockchain

type MessageService interface {
	RequestBlock() error
	ResponseBlock() error
}
