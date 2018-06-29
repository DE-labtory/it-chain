package p2p

type PeerService struct {
	leaderRepository LeaderRepository
	peerRepository PeerRepository
}

func (service *PeerService) GetLeader() (Leader){
	leader := service.leaderRepository.GetLeader()
	return leader
}

func (service *PeerService) GetPeerList() []Peer{
	peerList, _ := service.peerRepository.FindAll()
	return peerList
}

//get peer leader table
func (service *PeerService) GetPeerLeaderTable() PeerLeaderTable{
	leader := service.leaderRepository.GetLeader()
	peerList, _ := service.peerRepository.FindAll()
	peerLeaderTable := PeerLeaderTable{
		Leader:leader,
		PeerList:peerList,
	}
	return peerLeaderTable
}
