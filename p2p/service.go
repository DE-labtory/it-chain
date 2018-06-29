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

func (service *Service) GetPeerTable() PeerTable{
	leader := service.leaderRepository.GetLeader()
	peerList, _ := service.peerRepository.FindAll()
	peerTable := PeerTable{
		Leader:leader,
		PeerList:peerList,
	}
	return peerTable
}

