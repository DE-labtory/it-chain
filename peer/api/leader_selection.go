package api

import (
	"github.com/it-chain/it-chain-Engine/conf"
	"github.com/it-chain/it-chain-Engine/peer/domain/model"
	"github.com/it-chain/it-chain-Engine/peer/domain/repository"
	"github.com/it-chain/it-chain-Engine/peer/domain/service"
)

type LeaderSelection struct {
	peerTableService *service.PeerTable
	messageProducer  service.MessageProducer
	peerRepository   repository.Peer
	myInfo           *model.Peer
}

func NewLeaderSelectionApi(repo repository.Peer, messageProducer service.MessageProducer, myInfo *model.Peer) (*LeaderSelection, error) {
	LeaderSelectionApi := &LeaderSelection{
		peerTableService: service.NewPeerTableService(repo, myInfo),
		messageProducer:  messageProducer,
		peerRepository:   repo,
		myInfo:           myInfo,
	}
	bootNodeIp := conf.GetConfiguration().Common.BootNodeIp
	myIp := conf.GetConfiguration().Common.NodeIp
	if bootNodeIp == myIp {
		err := LeaderSelectionApi.peerTableService.SetLeader(LeaderSelectionApi.myInfo)
		return nil, err
	} else {
		err := LeaderSelectionApi.messageProducer.RequestLeaderInfo(bootNodeIp)
		return nil, err
	}
	return LeaderSelectionApi, nil
}

func (ls *LeaderSelection) RequestChangeLeader() error {
	panic("implement please")
}

func (ls *LeaderSelection) ChangeLeader(peer *model.Peer) error {
	return ls.peerTableService.SetLeader(peer)
}
