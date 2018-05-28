package messaging

import "github.com/it-chain/it-chain-Engine/p2p"

type GrpcMessageHandler struct {
	nodeRepository   p2p.NodeRepository
	leaderRepository p2p.LeaderRepository
	dispatcher       *Dispatcher
}

func NewGrpcMessageHandler(nodeRepo *p2p.NodeRepository, leaderRepo *p2p.LeaderRepository, dispatcher *Dispatcher) *GrpcMessageHandler {
	return &GrpcMessageHandler{
		nodeRepository:   nodeRepo,
		leaderRepository: leaderRepo,
		dispatcher:       dispatcher,
	}
}

//todo implement
func (gmh *GrpcMessageHandler) HandleMessageReceive(command p2p.GrpcRequestCommand) {
	panic("need to implement")
}
