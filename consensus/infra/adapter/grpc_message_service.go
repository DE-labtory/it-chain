package adapter

import (
	"github.com/it-chain/engine/api_gateway"
	"github.com/it-chain/engine/consensus"
)

type GrpcMessageService struct {
	pQuery api_gateway.PeerQueryApi
}
func NewGrpcService(peerRepository *api_gateway.PeerRepository) *ParliamentService {
	return &ParliamentService{
		pQuery: api_gateway.NewPeerQueryApi(peerRepository),
	}
}

func (service GrpcMessageService) sendPrePrepareMsgToPeer(msg consensus.PrePrepareMsg) error{

	//peerList, err := service.pQuery.GetPeerList()
	return nil
}

func (service GrpcMessageService) sendPrepareMsgToPeer(msg consensus.PrepareMsg) error{
	return nil
}


