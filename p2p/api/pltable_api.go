package api

import "github.com/it-chain/it-chain-Engine/p2p"

type PLTableApi struct {
	leaderService p2p.LeaderService
	peerService p2p.PeerService
}


func (plta *PLTableApi) GetPLTable() p2p.PLTable {

	leader := plta.leaderService.Get()
	peerList, _ := plta.peerService.FindAll()

	peerLeaderTable := p2p.PLTable{
		Leader:   leader,
		PeerList: peerList,
	}

	return peerLeaderTable
}
