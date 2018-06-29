package p2p

type Service struct {
	leaderRepository LeaderRepository
	peerRepository PeerRepository
}

func (service *Service) GetLeader() (Leader){
	leader := service.leaderRepository.GetLeader()
	return leader
}

func (service *Service) GetPeerList() []Peer{
	peerList, _ := service.peerRepository.FindAll()
	return peerList
}

//get peer leader table
func (service *Service) GetPeerLeaderTable() PeerLeaderTable{
	leader := service.leaderRepository.GetLeader()
	peerList, _ := service.peerRepository.FindAll()
	peerLeaderTable := PeerLeaderTable{
		Leader:leader,
		PeerList:peerList,
	}
	return peerLeaderTable
}

