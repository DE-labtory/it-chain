package blockchain

type GrpcCommandService interface {
	RequestBlock() error
	ResponseBlock() error
}
