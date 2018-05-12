package service

import (
	"github.com/it-chain/it-chain-Engine/peer/domain/model"
	"github.com/it-chain/it-chain-Engine/peer/domain/repository"
)

type PeerTable struct {
	leader     *model.Peer
	myInfo     *model.Peer
	repository repository.Peer
}

func NewPeerTableService(peerRepo repository.Peer, myinfo *model.Peer) *PeerTable {
	peerRepo.Save(*myinfo)
	return &PeerTable{
		leader:     nil,
		myInfo:     myinfo,
		repository: peerRepo,
	}
}

func (pts *PeerTable) SetLeader(peer *model.Peer) error {
	find, err := pts.repository.FindById(peer.Id)
	if err != nil {
		return err
	}
	pts.leader = find
	return nil
}

func (pts *PeerTable) GetLeader() *model.Peer {
	return pts.leader
}
